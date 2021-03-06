package gameprocessor

import (
	"fmt"
	"strings"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executePlayerCommand(cmd command, player *switchers.Player) {
	if cmd.text == commandPause {
		if err := gp.PlayerRepository.SetPaused(player, true); err != nil {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			return
		}
	}

	if player.Paused {
		if cmd.text != commandResume {
			gp.Bot.SendMessage(player.ID, responseGamePaused)
			return
		}
		if err := gp.PlayerRepository.SetPaused(player, false); err != nil {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			return
		}
		gp.Bot.SendMessage(player.ID, responseGameResumed)
		return
	}

	switch player.State {
	case switchers.PlayerStateNew:
		if err := gp.PlayerRepository.SetState(player, switchers.PlayerStateAskName); err != nil {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			return
		}
		gp.Bot.SendMessage(player.ID, responseAskName)
		return

	case switchers.PlayerStateAskName:
		if err := gp.PlayerRepository.SetState(player, switchers.PlayerStateIdle); err != nil {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			return
		}
		if err := gp.PlayerRepository.SetName(player, cmd.text); err != nil {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			return
		}
		gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseNiceToMeet, player.Name))
		return

	case switchers.PlayerStateIdle:
		if cmd.text == commandSetName {
			if err := gp.PlayerRepository.SetState(player, switchers.PlayerStateAskName); err != nil {
				gp.Bot.SendMessage(player.ID, responseSomethingWrong)
				return
			}
			gp.Bot.SendMessage(player.ID, responseSetName)
			return
		}

		if cmd.text == commandLeaders {
			leaders, err := gp.PlayerRepository.GetTop(5)
			if err != nil {
				gp.Bot.SendMessage(player.ID, responseSomethingWrong)
				return
			}
			var response string
			i := 1
			for _, leader := range leaders {
				response += fmt.Sprintf("%d. %s — %d\n", i, leader.Name, leader.Score)
				i++
			}
			gp.Bot.SendMessage(player.ID, fmt.Sprintf(responseLeaders, response, player.Score))
			return
		}

	case switchers.PlayerStateInGame:
		round, err := gp.RoundRepository.GetActiveRound()
		if err != nil || round.ID == "" {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			gp.Logger.Log("msg", "player is in ingame state, but no active round exists", "playerid", player.ID)
			return
		}
		playersTeamIndex := -1
		for i, team := range round.Teams {
			var exists bool
			if team.State == teamStateGathering {
				_, exists = team.GatheringPlayers[player.ID]
			} else {
				_, exists = team.ActualPlayers[player.ID]
			}
			if exists {
				playersTeamIndex = i
				break
			}
		}
		if playersTeamIndex == -1 {
			gp.Bot.SendMessage(player.ID, responseSomethingWrong)
			gp.Logger.Log("msg", "player is in ingame state, but no team found", "playerid", player.ID)
			return
		}

		switch round.Teams[playersTeamIndex].State {
		case teamStateGathering:
			if strings.ToLower(cmd.text) == commandGathered {
				if err = gp.RoundRepository.AddTeamMemberToActual(round, playersTeamIndex, player.ID); err != nil {
					gp.Bot.SendMessage(player.ID, responseSomethingWrong)
					gp.Logger.Log("msg", "failed to set player gathered state", "playerid", player.ID, "teamindex", playersTeamIndex)
					return
				}
				gp.Bot.SendMessage(player.ID, responsePlayerGathered)
				return
			}
			gp.Bot.SendMessage(player.ID, responseGatheringInstructions)
			return

		case teamStatePlaying:
			if strings.ToLower(cmd.text) == commandGathered {
				gp.Bot.SendMessage(player.ID, responseGatherNotAnswer)
				return
			}

			answer := switchers.Answer{MessageID: cmd.id, Text: cmd.text, OwnerID: player.ID}
			if err = gp.RoundRepository.SetTeamAnswer(round, playersTeamIndex, &answer); err != nil {
				gp.Bot.SendMessage(player.ID, responseSomethingWrong)
				gp.Logger.Log("msg", "failed to set team answer", "playerid", player.ID, "teamindex", playersTeamIndex, "answer", cmd.text)
				return
			}
			gp.Bot.SendMessage(player.ID, responsePlayerAnswered)
			return

		case teamStateModeration:
			gp.Bot.SendMessage(player.ID, responseWaitForModeration)
			return
		}
	}

	gp.Bot.SendMessage(player.ID, responseDefault)
}
