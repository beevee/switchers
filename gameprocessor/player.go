package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

const (
	stateNew     = ""
	stateAskName = "askname"
	stateIdle    = "idle"

	commandNewRound = "/newround"
	commandResign   = "/resign"

	commandSetName = "/setname"
	commandPause   = "/pause"
	commandResume  = "/resume"
)

// ExecuteCommand takes text command from a player, updates internal state and returns response
func (gp *GameProcessor) ExecuteCommand(command string, player *switchers.Player) string {
	var response string

	if player.Trump {
		response = gp.executeTrumpCommand(command, player)
	} else {
		response = gp.executePlayerCommand(command, player)
	}

	if err := gp.PlayerRepository.SavePlayer(player); err != nil {
		return "Что-то пошло не так, попробуй еще раз."
	}
	return response
}

func (gp *GameProcessor) executePlayerCommand(command string, player *switchers.Player) string {
	if command == commandPause {
		player.Paused = true
	}
	if player.Paused {
		if command != commandResume {
			return fmt.Sprintf("Участие в игре приостановлено. Ничего не сможешь делать, пока не напишешь %s.", commandResume)
		}
		player.Paused = false
		return "Участие в игре возобновлено. Продолжай как ни в чем не бывало."
	}

	switch player.State {
	case stateNew:
		player.State = stateAskName
		return "Привет! Чтобы стать участником Свитчеров, напиши в ответ свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."

	case stateAskName:
		player.State = stateIdle
		player.Name = command
		return fmt.Sprintf("Приятно познакомиться, %s. Теперь жди инструкции. Они могут приходить в любой момент, так что держи телефон включенным! Чтобы приостановить участие в игре, напиши /pause.", player.Name)

	case stateIdle:
		if command == "/setname" {
			player.State = stateAskName
			return "Напиши свое имя. Важно, чтобы другие участники могли тебя узнать, так что не пиши ерунду."
		}
	}

	return fmt.Sprintf("Жди инструкции или напиши какую-нибудь команду. Я понимаю:\n\n%s — изменить имя\n%s — приостановить участие в игре", commandSetName, commandPause)
}
