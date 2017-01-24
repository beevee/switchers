package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

// ExecuteCommand takes text command from a player, updates internal state and returns response
func (gp *GameProcessor) ExecuteCommand(command string, playerID string) {
	player, created, err := gp.PlayerRepository.GetOrCreatePlayer(playerID)
	if err != nil {
		gp.Logger.Log("msg", "error retrieving player profile while executing command", "error", err)
	}
	if created {
		gp.Logger.Log("msg", "created player profile", "playerid", player.ID)
	}

	if command == gp.TrumpCode {
		player.Trump = true
	}

	if player.Trump {
		gp.executeTrumpCommand(command, player)
	} else {
		gp.executePlayerCommand(command, player)
	}

	if err := gp.PlayerRepository.SavePlayer(player); err != nil {
		gp.Logger.Log("msg", "error saving player profile after executing command", "playerid", player.ID)
		gp.Bot.SendMessage(player.ID, "Что-то пошло не так, возможно надо попробовать еще раз.")
	}
}

func (gp *GameProcessor) executePlayerCommand(command string, player *switchers.Player) {
	response := fmt.Sprintf("Жди инструкции или напиши какую-нибудь команду. Я понимаю:\n\n%s — изменить имя\n%s — приостановить участие в игре", commandSetName, commandPause)

	if command == commandPause {
		player.Paused = true
		player.State = playerStateIdle
	}

	if player.Paused {
		if command != commandResume {
			response = fmt.Sprintf("Участие в игре приостановлено. Ничего не сможешь делать, пока не напишешь %s.", commandResume)
		} else {
			player.Paused = false
			response = "Участие в игре возобновлено. Продолжай как ни в чем не бывало."
		}
	} else {
		switch player.State {
		case playerStateNew:
			player.State = playerStateAskName
			response = "Привет! Чтобы стать участником Свитчеров, напиши в ответ свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."

		case playerStateAskName:
			player.State = playerStateIdle
			player.Name = command
			response = fmt.Sprintf("Приятно познакомиться, %s. Теперь жди инструкции. Они могут приходить в любой момент, так что держи телефон включенным! Чтобы приостановить участие в игре, напиши /pause.", player.Name)

		case playerStateIdle:
			if command == "/setname" {
				player.State = playerStateAskName
				response = "Напиши свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."
			}
		}
	}

	gp.Bot.SendMessage(player.ID, response)
}
