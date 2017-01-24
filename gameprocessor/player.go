package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

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
