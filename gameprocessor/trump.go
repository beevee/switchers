package gameprocessor

import (
	"fmt"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(command string, player *switchers.Player) string {
	switch command {
	case commandNewRound:
		ar, err := gp.RoundRepository.GetActiveRound()
		if ar.ID != "" {
			return "Сейчас уже идет раунд, новый начать нельзя."
		}
		if err != nil {
			return fmt.Sprintf("Произошла ошибка при поиске активного раунда: %s", err)
		}

		round, err := gp.RoundRepository.CreateActiveRound()
		if err != nil {
			return fmt.Sprintf("Произошла ошибка при создании активного раунда: %s", err)
		}

		if err = gp.PopulateRound(round); err != nil {
			return fmt.Sprintf("Произошла ошибка при генерации активного раунда: %s", err)
		}

		if err = gp.RoundRepository.SaveRound(round); err != nil {
			return fmt.Sprintf("Произошла ошибка при сохранении сгенерированного активного раунда: %s", err)
		}

		return "Начался новый раунд."
	case commandResign:
		player.Trump = false
		return "Отставка принята."
	}

	return fmt.Sprintf("Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n%s — запустить новый раунд\n%s — подать в отставку", commandNewRound, commandResign)
}
