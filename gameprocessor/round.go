package gameprocessor

import (
	"errors"
	"math/rand"
	"time"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) populateRound(round *switchers.Round) error {
	players, err := gp.PlayerRepository.GetAllPlayers()
	if err != nil {
		return err
	}
	gp.Logger.Log("msg", "retrieved players for new round", "count", len(players))

	eligiblePlayerIDs := make([]string, 0, len(players))
	for _, player := range players {
		if !player.Trump && !player.Paused {
			eligiblePlayerIDs = append(eligiblePlayerIDs, player.ID)
		}
	}

	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle
	for i := range eligiblePlayerIDs {
		j := rand.Intn(i + 1)
		eligiblePlayerIDs[i], eligiblePlayerIDs[j] = eligiblePlayerIDs[j], eligiblePlayerIDs[i]
	}

	teamCount := len(eligiblePlayerIDs) / teamMinSize
	gp.Logger.Log("msg", "calculated team count", "count", teamCount, "minsize", teamMinSize)

	gatheringTasks, err := gp.TaskRepository.GetAllGatheringTasks()
	if err != nil {
		return err
	}
	if len(gatheringTasks) < teamCount {
		return errors.New("not enough gathering tasks for a full round")
	}
	gp.Logger.Log("msg", "retrieved gathering tasks for new round", "count", len(gatheringTasks))

	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle
	for i := range gatheringTasks {
		j := rand.Intn(i + 1)
		gatheringTasks[i], gatheringTasks[j] = gatheringTasks[j], gatheringTasks[i]
	}

	round.Teams = make([]switchers.Team, teamCount)
	for i := range round.Teams {
		round.Teams[i].State = teamStateGathering
		round.Teams[i].PlayerIDs = make(map[string]bool)
		round.Teams[i].GatheringTask = gatheringTasks[i]
		round.Teams[i].GatheringDeadline = round.StartTime.Add(time.Minute * time.Duration(gatheringTasks[i].TimeLimitMinutes))
	}
	for i, playerID := range eligiblePlayerIDs {
		teamNumber := i % teamCount
		round.Teams[teamNumber].PlayerIDs[playerID] = false
	}
	gp.Logger.Log("msg", "round population finished")

	return nil
}