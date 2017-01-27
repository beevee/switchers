package gameprocessor

import (
	"fmt"
	"time"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) gameProgressor() error {
	ticker := time.NewTicker(10 * time.Second)

	for {
		select {
		case <-ticker.C:
			round, err := gp.RoundRepository.GetActiveRound()
			if err != nil {
				gp.Logger.Log("msg", "failed to retrieve active round while enforcing deadlines", "error", err)
				break
			}
			if round.ID == "" {
				break
			}
			now := time.Now()
			roundDeactivationRequired := true

		TEAMLOOP:
			for i, team := range round.Teams {
				if team.State == teamStateGathering || team.State == teamStatePlaying || team.State == teamStateModeration {
					roundDeactivationRequired = false
				}
				if team.State == teamStateGathering {
					if len(team.ActualPlayers) >= gp.TeamQuorum {
						if team.MissingPlayersDeadline.IsZero() {
							if err = gp.RoundRepository.SetTeamMissingPlayersDeadline(round, i, now.Add(1*time.Minute)); err != nil {
								gp.Logger.Log("msg", "failed to set gathered team actual deadline", "teamindex", i, "error", err)
							}
							continue
						}
						if now.Before(team.MissingPlayersDeadline) {
							continue
						}
						if err = gp.RoundRepository.SetTeamActualDeadline(round, i, now.Add(time.Duration(team.ActualTask.TimeLimitMinutes)*time.Minute)); err != nil {
							gp.Logger.Log("msg", "failed to set gathered team actual deadline", "teamindex", i, "error", err)
							continue
						}
						if err = gp.RoundRepository.SetTeamState(round, i, teamStatePlaying); err != nil {
							gp.Logger.Log("msg", "failed to set playing team state", "teamindex", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "team quorum gathered", "teamindex", i, "count", len(team.ActualPlayers))
						gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamGotQuorum, i))
						gp.notifyActualTeamMembers(team, team.ActualTask.Text+responseActualTaskSuffix)
						for playerID := range team.GatheringPlayers {
							if _, exists := team.ActualPlayers[playerID]; !exists {
								if err = gp.RoundRepository.AddTeamMemberToMissing(round, i, playerID); err != nil {
									gp.Logger.Log("msg", "failed to add player to missing", "playerid", playerID, "teamindex", i, "error", err)
								}
								gp.PlayerRepository.SetState(&switchers.Player{ID: playerID}, playerStateIdle)
								gp.Bot.SendMessage(playerID, responsePlayerFailedToGather)
							}
						}
						continue
					}

					if now.After(team.GatheringDeadline) {
						if err = gp.RoundRepository.SetTeamState(round, i, teamStateLost); err != nil {
							gp.Logger.Log("msg", "failed to set timeouted team state (gathering)", "teamindex", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "timeouted a team (gathering)", "teamindex", i)
						gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamFailedToGather, i))
						gp.updateGatheringTeamMemberStates(team, playerStateIdle)
						gp.notifyGatheringTeamMembers(team, responseTeamFailedToGather)
					}
				}

				if team.State == teamStatePlaying {
					if team.Answer != "" {
						for _, correctAnswer := range team.ActualTask.CorrectAnswers {
							if team.Answer == correctAnswer {
								if err = gp.RoundRepository.SetTeamState(round, i, teamStateWon); err != nil {
									gp.Logger.Log("msg", "failed to set won team state (playing)", "teamindex", i, "error", err)
									continue TEAMLOOP
								}
								gp.Logger.Log("msg", "team won by answering correctly", "teamindex", i, "answer", team.Answer)
								gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamWon, i))
								gp.updateActualTeamMemberStates(team, playerStateIdle)
								gp.increaseActualTeamMemberScores(team)
								gp.notifyActualTeamMembers(team, responseTeamWon)
								continue TEAMLOOP
							}
						}

						if err = gp.RoundRepository.SetTeamState(round, i, teamStateModeration); err != nil {
							gp.Logger.Log("msg", "failed to set moderation team state (playing)", "teamindex", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "team answer is available for moderation", "teamindex", i, "answer", team.Answer)
						gp.notifyTrumps(fmt.Sprintf(responseTrumpModerationRequired, i))
						gp.notifyActualTeamMembers(team, responseModerationRequired)
						continue
					}

					if now.After(team.ActualDeadline) {
						if err = gp.RoundRepository.SetTeamState(round, i, teamStateLost); err != nil {
							gp.Logger.Log("msg", "failed to timeouted team state (playing)", "teamindex", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "timeouted a team (playing)", "teamindex", i)
						gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamFailedToAnswer, i))
						gp.updateActualTeamMemberStates(team, playerStateIdle)
						gp.notifyActualTeamMembers(team, responseTeamFailedToAnswer)
					}
				}
			}

			if roundDeactivationRequired {
				if err = gp.RoundRepository.DeactivateRound(round); err != nil {
					gp.Logger.Log("msg", "failed to deactivate round", "error", err)
					break
				}
				gp.Logger.Log("msg", "deactivated a round", "roundid", round.ID)
				gp.notifyTrumps(responseTrumpActiveRoundFinished)
			}
		case <-gp.tomb.Dying():
			return nil
		}
	}
}
