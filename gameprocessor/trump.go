package gameprocessor

import (
	"fmt"
	"strings"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) executeTrumpCommand(command string, player *switchers.Player) {
	if player.State == playerStateModerating {
		round, err := gp.RoundRepository.GetActiveRound()
		if round.ID == "" {
			gp.Bot.SendMessage(player.ID, "Сейчас нет активного раунда, модерировать нечего.")
			return
		}
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при поиске активного раунда: %s", err))
			return
		}

		team := round.Teams[player.ModeratingTeamIndex]
		if team.State != teamStateModeration {
			gp.Bot.SendMessage(player.ID, "Эту команду кто-то уже отмодерировал.")
		} else {
			switch strings.ToLower(command) {
			case "да":
				if err = gp.RoundRepository.SetTeamState(round, player.ModeratingTeamIndex, teamStateWon); err != nil {
					gp.Logger.Log("msg", "failed to set won team state (playing)", "index", player.ModeratingTeamIndex, "error", err)
				}
				gp.Logger.Log("msg", "team won by answering correctly after moderation", "index", player.ModeratingTeamIndex, "answer", team.Answer)
				gp.notifyTrumps(fmt.Sprintf("Команда %d выиграла, дав правильный ответ.", player.ModeratingTeamIndex))
				gp.updateActualTeamMemberStates(team, playerStateIdle)
				gp.increaseActualTeamMemberScores(team)
				gp.notifyActualTeamMembers(team, "Вы победили и получаете кучу очков! Ждите следующий раунд.")

			case "нет":
				if err = gp.RoundRepository.SetTeamState(round, player.ModeratingTeamIndex, teamStateLost); err != nil {
					gp.Logger.Log("msg", "failed to set lost team state (playing)", "index", player.ModeratingTeamIndex, "error", err)
				}
				gp.Logger.Log("msg", "team lost by answering incorrectly after moderation (playing)", "index", player.ModeratingTeamIndex)
				gp.notifyTrumps(fmt.Sprintf("Команда %d проиграла, дав неправильный ответ.", player.ModeratingTeamIndex))
				gp.updateActualTeamMemberStates(team, playerStateIdle)
				gp.notifyActualTeamMembers(team, "Вы проиграли, потому что ответили неправильно. Теперь вы не получите кучу очков :(")
			}
		}

		gp.Bot.SendMessage(player.ID, "We made America great again!")
		if err := gp.PlayerRepository.SetState(player, playerStateIdle); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при выходе из режима модерации: %s", err))
		}
		return
	}

	switch command {
	case commandNewRound:
		if err := gp.startNewRound(); err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при генерации нового раунда: %s", err))
			return
		}
		gp.Bot.SendMessage(player.ID, "Начался новый раунд.")
		return

	case commandModerate:
		round, err := gp.RoundRepository.GetActiveRound()
		if round.ID == "" {
			gp.Bot.SendMessage(player.ID, "Сейчас нет активного раунда, модерировать нечего.")
			return
		}
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при поиске активного раунда: %s", err))
			return
		}
		for i, team := range round.Teams {
			if team.State == teamStateModeration {
				if err := gp.PlayerRepository.SetModeratingTeamIndex(player, i); err != nil {
					gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при поиске команды для модерации: %s", err))
					return
				}
				if err := gp.PlayerRepository.SetState(player, playerStateModerating); err != nil {
					gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при входе в режим модерации: %s", err))
					return
				}
				gp.Bot.SendMessage(player.ID, "Задание: "+team.ActualTask.Text)
				gp.Bot.SendMessage(player.ID, "Ответ: "+team.Answer)
				gp.Bot.SendMessage(player.ID, "Напиши \"да\", если ответ правильный. Напиши \"нет\", если ответ неправильный. Напиши что угодно другое, чтобы бросить это занятие и пойти строить стену на границе с Мексикой.")
				return
			}
		}

	case commandLeaders:
		leaders, err := gp.PlayerRepository.GetTop(500)
		if err != nil {
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка при чтении списка игроков: %s", err))
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
			gp.Bot.SendMessage(player.ID, fmt.Sprintf("Произошла ошибка в процедуре отставки: %s", err))
			return
		}
		gp.Bot.SendMessage(player.ID, "Отставка принята.")
		return
	}

	gp.Bot.SendMessage(player.ID, fmt.Sprintf("Добро пожаловать, господин президент. Издайте какой-нибудь указ:\n\n%s — запустить новый раунд\n%s — модерировать что-нибудь\n%s — посмотреть очки игроков\n%s — подать в отставку", commandNewRound, commandModerate, commandLeaders, commandResign))
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
