package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(command string, player *switchers.Player) {
	switch command {
	case commandNewRound:
		ar, err := gp.RoundRepository.GetActiveRound()
		if ar.ID != "" {
			gp.Bot.SendMessage(player.ID, "Сейчас уже идет раунд, новый начать нельзя.")
			return
		}
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при поиске активного раунда: %s", err))
			return
		}

		round, err := gp.RoundRepository.CreateActiveRound()
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при создании активного раунда: %s", err))
			return
		}

		if err = gp.populateRound(round); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при генерации активного раунда: %s", err))
			return
		}

		if err = gp.RoundRepository.SaveRound(round); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при сохранении сгенерированного активного раунда: %s", err))
			return
		}

		for _, team := range round.Teams {
			gp.notifyTeam(team, team.GatheringTask.Text)
			gp.updateTeamMemberStates(team, playerStateInGame)
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
