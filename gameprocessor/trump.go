package gameprocessor

import (
	"fmt"
	"strings"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(cmd command, player *switchers.Player) {
	if player.State == playerStateModerating {
		round, err := gp.RoundRepository.GetActiveRound()
		if round.ID == "" {
			gp.Bot.SendMessage(player.ID, responseTrumpNothingToModerate)
			return
		}
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
			return
		}

		team := round.Teams[player.ModeratingTeamIndex]
		if team.State != teamStateModeration {
			gp.Bot.SendMessage(player.ID, responseTrumpAlreadyModerated)
		} else {
			switch strings.ToLower(cmd.text) {
			case commandYes:
				if err = gp.RoundRepository.SetTeamState(round, player.ModeratingTeamIndex, teamStateWon); err != nil {
					gp.Logger.Log("msg", "failed to set won team state (playing)", "teamindex", player.ModeratingTeamIndex, "error", err)
				}
				gp.Logger.Log("msg", "team won by answering correctly after moderation", "teamindex", player.ModeratingTeamIndex, "answer", team.Answer)
				gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamWon, player.ModeratingTeamIndex))
				gp.updateActualTeamMemberStates(team, playerStateIdle)
				gp.increaseActualTeamMemberScores(team)
				gp.notifyActualTeamMembers(team, responseTeamWon)

			case commandNo:
				if err = gp.RoundRepository.SetTeamState(round, player.ModeratingTeamIndex, teamStateLost); err != nil {
					gp.Logger.Log("msg", "failed to set lost team state (playing)", "teamindex", player.ModeratingTeamIndex, "error", err)
				}
				gp.Logger.Log("msg", "team lost by answering incorrectly after moderation (playing)", "teamindex", player.ModeratingTeamIndex)
				gp.notifyTrumps(fmt.Sprintf(responseTrumpTeamLost, player.ModeratingTeamIndex))
				gp.updateActualTeamMemberStates(team, playerStateIdle)
				gp.notifyActualTeamMembers(team, responseTeamLost)
			}
		}

		gp.Bot.SendMessage(player.ID, "We made America great again!")
		if err := gp.PlayerRepository.SetState(player, playerStateIdle); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
		}
		return
	}

	switch cmd.text {
	case commandNewRound:
		if err := gp.startNewRound(); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
			return
		}
		gp.Bot.SendMessage(player.ID, responseTrumpRoundStarted)
		return

	case commandModerate:
		round, err := gp.RoundRepository.GetActiveRound()
		if round.ID == "" {
			gp.Bot.SendMessage(player.ID, responseTrumpNothingToModerate)
			return
		}
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
			return
		}
		for i, team := range round.Teams {
			if team.State == teamStateModeration {
				if err := gp.PlayerRepository.SetModeratingTeamIndex(player, i); err != nil {
					gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
					return
				}
				if err := gp.PlayerRepository.SetState(player, playerStateModerating); err != nil {
					gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
					return
				}
				gp.Bot.SendMessage(player.ID, responseTrumpTaskPrefix+team.ActualTask.Text)
				gp.Bot.ForwardMessage(player.ID, team.Answer.Text, team.Answer.MessageID, team.Answer.OwnerID)
				gp.Bot.SendMessage(player.ID, responseTrumpModerationInstructions)
				return
			}
		}

	case commandLeaders:
		leaders, err := gp.PlayerRepository.GetTop(500)
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
			return
		}
		var response string
		i := 1
		for _, leader := range leaders {
			response += fmt.Sprintf("%d. %s — %d\n", i, leader.Name, leader.Score)
			i++
		}
		gp.Bot.SendMessage(player.ID, response)
		return

	case commandResign:
		if err := gp.PlayerRepository.SetTrump(player, false); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseTrumpSomethingWrong, err))
			return
		}
		gp.Bot.SendMessage(player.ID, responseTrumpResigned)
		return
	}

	gp.Bot.SendMessage(player.ID, responseTrumpDefault)
}

func (gp *GameProcessor) notifyTrumps(message string) {
	trumpIDs, err := gp.PlayerRepository.GetAllTrumps()
	if err != nil {
		gp.Logger.Log("msg", "failed to notify Trumps", "message", message, "error", err)
	} else {
		for trumpID := range trumpIDs {
			gp.Bot.SendMessage(trumpID, message)
		}
	}
}
