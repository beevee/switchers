package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(command string, player *switchers.Player) {
	switch command {
	case commandNewRound:
		if err := gp.startNewRound(); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при генерации нового раунда: %s", err))
			return
		}
		gp.Bot.SendMessage(player.ID, "Начался новый раунд.")
		return

	case commandResign:
		if err := gp.PlayerRepository.SetTrump(player, false); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка в процедуре отставки: %s", err))
			return
		}
		gp.Bot.SendMessage(player.ID, "Отставка принята.")
		return
	}

	gp.Bot.SendMessage(player.ID, fmt.Sprintf("Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n%s — запустить новый раунд\n%s — подать в отставку", commandNewRound, commandResign))
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
