package gameprocessor

import "github.com/beevee/switchers"

func (gp *GameProcessor) notifyGatheringTeamMembers(team *switchers.Team, message string) {
	for playerID := range team.GatheringPlayers {
		gp.Bot.SendMessage(playerID, message)
	}
}

func (gp *GameProcessor) notifyActualTeamMembers(team *switchers.Team, message string) {
	for playerID := range team.ActualPlayers {
		gp.Bot.SendMessage(playerID, message)
	}
}

func (gp *GameProcessor) updateGatheringTeamMemberStates(team *switchers.Team, state string) {
	for playerID := range team.GatheringPlayers {
		if err := gp.PlayerRepository.SetState(&switchers.Player{ID: playerID}, state); err != nil {
			gp.Logger.Log("msg", "failed to update team member state", "state", state, "playerid", playerID, "error", err)
		}
	}
}

func (gp *GameProcessor) updateActualTeamMemberStates(team *switchers.Team, state string) {
	for playerID := range team.ActualPlayers {
		if err := gp.PlayerRepository.SetState(&switchers.Player{ID: playerID}, state); err != nil {
			gp.Logger.Log("msg", "failed to update team member state", "state", state, "playerid", playerID, "error", err)
		}
	}
}

func (gp *GameProcessor) increaseActualTeamMemberScores(team *switchers.Team) {
	for playerID := range team.ActualPlayers {
		player, _, err := gp.PlayerRepository.GetOrCreatePlayer(playerID)
		if err != nil {
			gp.Logger.Log("msg", "failed to increase team member score", "playerid", playerID, "error", err)
			continue
		}
		if err := gp.PlayerRepository.IncreaseScore(player); err != nil {
			gp.Logger.Log("msg", "failed to increase team member score", "playerid", playerID, "error", err)
		}
	}
}
