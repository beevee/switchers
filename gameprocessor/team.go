package gameprocessor

import "github.com/beevee/switchers"

func (gp *GameProcessor) notifyTeam(team switchers.Team, message string) {
	for playerID := range team.PlayerIDs {
		gp.Bot.SendMessage(playerID, message)
	}
}

func (gp *GameProcessor) updateTeamMemberStates(team switchers.Team, state string) {
	for playerID := range team.PlayerIDs {
		if err := gp.PlayerRepository.SetPlayerState(playerID, state); err != nil {
			gp.Logger.Log("msg", "failed to update team member state", "state", state, "playerid", playerID, "error", err)
		}
	}
}
