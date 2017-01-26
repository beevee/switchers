package gameprocessor

import (
	"fmt"
	"time"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) roundDeactivator() error {
	ticker := time.NewTicker(10 * time.Second)

	for {
	SELECT:
		select {
		case <-ticker.C:
			round, err := gp.RoundRepository.GetActiveRound()
			if err != nil {
				gp.Logger.Log("msg", "failed to retrieve active round while checking for deactivation", "error", err)
				break
			}
			if round.ID == "" {
				break
			}
			for _, team := range round.Teams {
				if team.State == teamStateGathering || team.State == teamStatePlaying || team.State == teamStateModeration {
					break SELECT
				}
			}
			if err = gp.RoundRepository.DeactivateRound(round); err != nil {
				gp.Logger.Log("msg", "failed to deactivate round", "error", err)
				break
			}
			gp.Logger.Log("msg", "deactivated a round", "newid", round.ID)
			gp.notifyTrumps("Активный раунд завершился.")
		case <-gp.tomb.Dying():
			return nil
		}
	}
}

func (gp *GameProcessor) deadlineEnforcer() error {
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

		TEAMLOOP:
			for i, team := range round.Teams {
				if team.State == teamStateGathering {
					if len(team.ActualPlayers) >= gp.TeamQuorum {
						if team.MissingPlayersDeadline.IsZero() {
							if err = gp.RoundRepository.SetTeamMissingPlayersDeadline(round, i, now.Add(1*time.Minute)); err != nil {
								gp.Logger.Log("msg", "failed to set gathered team actual deadline", "index", i, "error", err)
							}
							continue
						}
						if now.Before(team.MissingPlayersDeadline) {
							continue
						}
						if err = gp.RoundRepository.SetTeamActualDeadline(round, i, now.Add(time.Duration(team.ActualTask.TimeLimitMinutes)*time.Minute)); err != nil {
							gp.Logger.Log("msg", "failed to set gathered team actual deadline", "index", i, "error", err)
							continue
						}
						if err = gp.RoundRepository.SetTeamState(round, i, teamStatePlaying); err != nil {
							gp.Logger.Log("msg", "failed to set playing team state", "index", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "team quorum gathered", "index", i, "count", len(team.ActualPlayers))
						gp.notifyTrumps(fmt.Sprintf("Команда %d набрала кворум, даем задачу.", i))
						gp.notifyActualTeamMembers(team, team.ActualTask.Text)
						for playerID := range team.GatheringPlayers {
							if _, exists := team.ActualPlayers[playerID]; !exists {
								if err = gp.RoundRepository.AddTeamMemberToMissing(round, i, playerID); err != nil {
									gp.Logger.Log("msg", "failed to add player to missing", "playerid", playerID, "index", i, "error", err)
								}
								gp.PlayerRepository.SetState(&switchers.Player{ID: playerID}, playerStateIdle)
								gp.Bot.SendMessage(playerID, "Нужно было вовремя написать \"тут\", а у тебя не получилось. Жди теперь следующий раунд.")
							}
						}
						continue
					}

					if now.After(team.GatheringDeadline) {
						if err = gp.RoundRepository.SetTeamState(round, i, teamStateLost); err != nil {
							gp.Logger.Log("msg", "failed to set timeouted team state (gathering)", "index", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "timeouted a team (gathering)", "index", i)
						gp.notifyTrumps(fmt.Sprintf("У команды %d закончилось время на сборы, они проиграли.", i))
						gp.notifyGatheringTeamMembers(team, "Время вышло :( Этот раунд вы проиграли, потому что не собрали команду вовремя. Но в следующий раз повезет! Ждите следующий раунд.")
						gp.updateGatheringTeamMemberStates(team, playerStateIdle)
					}
				}

				if team.State == teamStatePlaying {
					if team.Answer != "" {
						for _, correctAnswer := range team.ActualTask.CorrectAnswers {
							if team.Answer == correctAnswer {
								if err = gp.RoundRepository.SetTeamState(round, i, teamStateWon); err != nil {
									gp.Logger.Log("msg", "failed to set won team state (playing)", "index", i, "error", err)
									continue TEAMLOOP
								}
								gp.Logger.Log("msg", "team won by answering correctly", "index", i, "answer", team.Answer)
								gp.notifyTrumps(fmt.Sprintf("Команда %d выиграла, дав правильный ответ.", i))
								gp.notifyActualTeamMembers(team, "Вы победили и получаете кучу очков! Ждите следующий раунд.")
								gp.updateActualTeamMemberStates(team, playerStateIdle)
								gp.increaseActualTeamMemberScores(team)
								continue TEAMLOOP
							}
						}

						if err = gp.RoundRepository.SetTeamState(round, i, teamStateModeration); err != nil {
							gp.Logger.Log("msg", "failed to set moderation team state (playing)", "index", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "team answer is available for moderation", "index", i, "answer", team.Answer)
						gp.notifyTrumps(fmt.Sprintf("Команда %d дала ответ, требуется модерация.", i))
						gp.notifyActualTeamMembers(team, "Ваш ответ направлен на модерацию, ждите решения.")
						continue
					}

					if now.After(team.ActualDeadline) {
						if err = gp.RoundRepository.SetTeamState(round, i, teamStateLost); err != nil {
							gp.Logger.Log("msg", "failed to timeouted team state (playing)", "index", i, "error", err)
							continue
						}
						gp.Logger.Log("msg", "timeouted a team (playing)", "index", i)
						gp.notifyTrumps(fmt.Sprintf("У команды %d закончилось время на ответ, они проиграли.", i))
						gp.notifyActualTeamMembers(team, "Время вышло :( Этот раунд вы проиграли, потому что не ответили на задачу вовремя. Но в следующий раз повезет! Ждите следующий раунд.")
						gp.updateActualTeamMemberStates(team, playerStateIdle)
					}
				}
			}
		case <-gp.tomb.Dying():
			return nil
		}
	}
}
