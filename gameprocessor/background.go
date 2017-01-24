package gameprocessor

import (
	"fmt"
	"time"
)

func (gp *GameProcessor) roundDeactivator() error {
	ticker := time.NewTicker(time.Minute)

	for {
	S:
		select {
		case <-ticker.C:
			round, err := gp.RoundRepository.GetActiveRound()
			if err != nil {
				gp.Logger.Log("msg", "failed to retrieve active round while checking for deactivation", "error", err)
				break
			}
			if round.ID == "" {
				break
			}
			for _, team := range round.Teams {
				if team.State == teamStateGathering || team.State == teamStatePlaying || team.State == teamStateModeration {
					break S
				}
			}
			if err = gp.RoundRepository.DeactivateRound(round); err != nil {
				gp.Logger.Log("msg", "failed to deactivate round", "error", err)
				break
			}
			gp.Logger.Log("msg", "deactivated a round", "newid", round.ID)
			gp.notifyTrumps("Активный раунд завершился.")
		case <-gp.tomb.Dying():
			return nil
		}
	}
}

func (gp *GameProcessor) deadlineEnforcer() error {
	ticker := time.NewTicker(time.Minute)

	for {
		select {
		case <-ticker.C:
			round, err := gp.RoundRepository.GetActiveRound()
			if err != nil {
				gp.Logger.Log("msg", "failed to retrieve active round while enforcing deadlines", "error", err)
				break
			}
			if round.ID == "" {
				break
			}
			now := time.Now()
			for i, team := range round.Teams {
				if team.State == teamStateGathering && now.After(team.GatheringDeadline) {
					team.State = teamStateLost
					if err = gp.RoundRepository.SaveTeam(round, i, team); err != nil {
						gp.Logger.Log("msg", "failed to save timeouted team (gathering)", "index", i, "error", err)
					} else {
						gp.Logger.Log("msg", "timeouted a team (gathering)", "index", i)
						gp.notifyTrumps(fmt.Sprintf("У команды %d закончилось время на сборы, они проиграли.", i))
					}
				}
			}
		case <-gp.tomb.Dying():
			return nil
		}
	}
}
