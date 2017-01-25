package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executePlayerCommand(command string, player *switchers.Player) {
	if command == commandPause {
		if err := gp.PlayerRepository.SetPaused(player, true); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			return
		}
	}

	if player.Paused {
		if command != commandResume {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Участие в игре приостановлено. Ничего не сможешь делать, пока не напишешь %s.", commandResume))
			return
		}
		if err := gp.PlayerRepository.SetPaused(player, false); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			return
		}
		gp.Bot.SendMessage(player.ID, "Участие в игре возобновлено. Продолжай как ни в чем не бывало.")
		return
	}

	switch player.State {
	case playerStateNew:
		if err := gp.PlayerRepository.SetState(player, playerStateAskName); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			return
		}
		gp.Bot.SendMessage(player.ID, "Привет! Чтобы стать участником Свитчеров, напиши в ответ свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду.")
		return

	case playerStateAskName:
		if err := gp.PlayerRepository.SetState(player, playerStateIdle); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			return
		}
		if err := gp.PlayerRepository.SetName(player, command); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			return
		}
		gp.Bot.SendMessage(player.ID, fmt.Sprintf("Приятно познакомиться, %s. Теперь жди инструкции. Они могут приходить в любой момент, так что держи телефон включенным! Чтобы приостановить участие в игре, напиши /pause.", player.Name))
		return

	case playerStateIdle:
		if command == "/setname" {
			if err := gp.PlayerRepository.SetState(player, playerStateAskName); err != nil {
				gp.Bot.SendMessage(player.ID, somethingWrongResponse)
				return
			}
			gp.Bot.SendMessage(player.ID, "Напиши свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду.")
			return
		}

	case playerStateInGame:
		round, err := gp.RoundRepository.GetActiveRound()
		if err != nil || round.ID == "" {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			gp.Logger.Log("msg", "player is in ingame state, but no active round exists", "playerid", player.ID)
			return
		}
		playersTeamIndex := -1
		for i, team := range round.Teams {
			_, exists := team.PlayerIDs[player.ID]
			if exists {
				playersTeamIndex = i
				break
			}
		}
		if playersTeamIndex == -1 {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			gp.Logger.Log("msg", "player is in ingame state, but no team found", "playerid", player.ID)
			return
		}

		if round.Teams[playersTeamIndex].State == teamStateGathering {
			if command == "тут" {
				if err = gp.RoundRepository.SetPlayerGathered(round, playersTeamIndex, player.ID); err != nil {
					gp.Bot.SendMessage(player.ID, somethingWrongResponse)
					gp.Logger.Log("msg", "failed to set player gathered state", "playerid", player.ID, "teamindex", playersTeamIndex)
					return
				}
				gp.Bot.SendMessage(player.ID, "Ждем отстающих еще немного и начинаем.")
				return
			}
			gp.Bot.SendMessage(player.ID, "Соберитесь в указанном месте. Как только соберетесь, каждый должен написать \"тут\".")
			return
		}

		if round.Teams[playersTeamIndex].State == teamStatePlaying {
			if err = gp.RoundRepository.SetTeamAnswer(round, playersTeamIndex, command); err != nil {
				gp.Bot.SendMessage(player.ID, somethingWrongResponse)
				gp.Logger.Log("msg", "failed to set team answer", "playerid", player.ID, "teamindex", playersTeamIndex, "answer", command)
				return
			}
			gp.Bot.SendMessage(player.ID, "Ответ принят.")
			return
		}
	}

	gp.Bot.SendMessage(player.ID, fmt.Sprintf("Жди инструкции или напиши какую-нибудь команду. Я понимаю:\n\n%s — изменить имя\n%s — приостановить участие в игре", commandSetName, commandPause))
}
