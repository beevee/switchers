package gameprocessor

import "github.com/beevee/switchers"

func (gp *GameProcessor) notifyTeam(team *switchers.Team, message string) {
	for playerID := range team.PlayerIDs {
		gp.Bot.SendMessage(playerID, message)
	}
}
