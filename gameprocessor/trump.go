package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(command string, player *switchers.Player) {
	response := fmt.Sprintf("Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n%s — запустить новый раунд\n%s — подать в отставку", commandNewRound, commandResign)

	switch command {
	case commandNewRound:
		ar, err := gp.RoundRepository.GetActiveRound()
		if ar.ID != "" {
			response = "Сейчас уже идет раунд, новый начать нельзя."
		}
		if err != nil {
			response = fmt.Sprintf("Произошла ошибка при поиске активного раунда: %s", err)
		}

		round, err := gp.RoundRepository.CreateActiveRound()
		if err != nil {
			response = fmt.Sprintf("Произошла ошибка при создании активного раунда: %s", err)
		}

		if err = gp.populateRound(round); err != nil {
			response = fmt.Sprintf("Произошла ошибка при генерации активного раунда: %s", err)
		}

		if err = gp.RoundRepository.SaveRound(round); err != nil {
			response = fmt.Sprintf("Произошла ошибка при сохранении сгенерированного активного раунда: %s", err)
		}

		for _, team := range round.Teams {
			gp.notifyTeam(team, team.GatheringTask.Text)
			gp.updateTeamMemberStates(team, playerStateInGame)
		}

		response = "Начался новый раунд."
	case commandResign:
		if err := gp.PlayerRepository.SetTrump(player, false); err != nil {
			response = fmt.Sprintf("Произошла ошибка в процедуре отставки: %s", err)
		} else {
			response = "Отставка принята."
		}
	}

	gp.Bot.SendMessage(player.ID, response)
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
