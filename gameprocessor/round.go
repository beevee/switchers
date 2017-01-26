package gameprocessor

import (
	"errors"
	"math/rand"
	"time"

	"github.com/beevee/switchers"
)

func (gp *GameProcessor) startNewRound() error {
	ar, err := gp.RoundRepository.GetActiveRound()
	if ar.ID != "" {
		return errors.New("cannot create new active round while another round is active")
	}
	if err != nil {
		return err
	}

	round, err := gp.generateRound()
	if err != nil {
		return err
	}

	if err = gp.RoundRepository.SaveActiveRound(round); err != nil {
		return err
	}

	for _, team := range round.Teams {
		gp.updateGatheringTeamMemberStates(team, playerStateInGame)
		gp.notifyGatheringTeamMembers(team, team.GatheringTask.Text)
	}

	return nil
}

func (gp *GameProcessor) generateRound() (*switchers.Round, error) {
	round := &switchers.Round{
		StartTime: time.Now(),
	}

	players, err := gp.PlayerRepository.GetAllPlayers()
	if err != nil {
		return nil, err
	}
	gp.Logger.Log("msg", "retrieved players for new round", "count", len(players))

	eligiblePlayers := make([]*switchers.Player, 0, len(players))
	for _, player := range players {
		if !player.Trump && !player.Paused {
			eligiblePlayers = append(eligiblePlayers, player)
		}
	}

	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle
	for i := range eligiblePlayers {
		j := rand.Intn(i + 1)
		eligiblePlayers[i], eligiblePlayers[j] = eligiblePlayers[j], eligiblePlayers[i]
	}

	teamCount := len(eligiblePlayers) / gp.TeamMinSize
	gp.Logger.Log("msg", "calculated team count", "count", teamCount, "minsize", gp.TeamMinSize)
	if teamCount == 0 {
		return nil, errors.New("not enough players to form a single team")
	}

	gatheringTasks, err := gp.TaskRepository.GetAllGatheringTasks()
	if err != nil {
		return nil, err
	}
	if len(gatheringTasks) < teamCount {
		return nil, errors.New("not enough gathering tasks for a full round")
	}
	gp.Logger.Log("msg", "retrieved gathering tasks for new round", "count", len(gatheringTasks))

	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle
	for i := range gatheringTasks {
		j := rand.Intn(i + 1)
		gatheringTasks[i], gatheringTasks[j] = gatheringTasks[j], gatheringTasks[i]
	}

	actualTasks, err := gp.TaskRepository.GetAllActualTasks()
	if err != nil {
		return nil, err
	}
	if len(actualTasks) < teamCount {
		return nil, errors.New("not enough actual tasks for a full round")
	}
	gp.Logger.Log("msg", "retrieved actual tasks for new round", "count", len(actualTasks))

	// https://en.wikipedia.org/wiki/Fisher–Yates_shuffle
	for i := range actualTasks {
		j := rand.Intn(i + 1)
		actualTasks[i], actualTasks[j] = actualTasks[j], actualTasks[i]
	}

	round.Teams = make([]*switchers.Team, teamCount)
	for i := range round.Teams {
		round.Teams[i] = &switchers.Team{
			State:             teamStateGathering,
			GatheringPlayers:  make(map[string]switchers.Player),
			GatheringTask:     gatheringTasks[i],
			ActualTask:        actualTasks[i],
			GatheringDeadline: round.StartTime.Add(time.Minute * time.Duration(gatheringTasks[i].TimeLimitMinutes)),
		}
	}
	for i, player := range eligiblePlayers {
		teamNumber := i % teamCount
		round.Teams[teamNumber].GatheringPlayers[player.ID] = *player
	}
	gp.Logger.Log("msg", "round population finished")

	return round, nil
}
