package firebase

import (
	"fmt"
	"strconv"
	"time"

	"github.com/beevee/switchers"
	"github.com/satori/go.uuid"
)

// RoundRepository persists round information in Firebase
type RoundRepository struct {
	Repository
}

// GetActiveRound looks up active round and returns it
func (rr *RoundRepository) GetActiveRound() (*switchers.Round, error) {
	ref, err := rr.firebase.Ref("rounds/active")
	if err != nil {
		return nil, err
	}

	round := &switchers.Round{}
	if err = ref.Value(round); err != nil {
		return nil, err
	}
	return round, nil
}

// DeactivateRound puts active round into archive
func (rr *RoundRepository) DeactivateRound(round *switchers.Round) error {
	if round.ID != "active" {
		return nil
	}

	round.ID = uuid.NewV4().String()
	ref, err := rr.firebase.Ref("rounds/" + round.ID)
	if err != nil {
		return err
	}
	if err = ref.Set(round); err != nil {
		return err
	}

	ref, err = rr.firebase.Ref("rounds/active")
	if err != nil {
		return err
	}

	return ref.Remove()
}

// SaveActiveRound saves active round
func (rr *RoundRepository) SaveActiveRound(round *switchers.Round) error {
	round.ID = "active"
	ref, err := rr.firebase.Ref("rounds/" + round.ID)
	if err != nil {
		return err
	}

	return ref.Set(round)
}

// AddTeamMemberToActual adds gathering player to actual team
func (rr *RoundRepository) AddTeamMemberToActual(round *switchers.Round, index int, playerID string) error {
	player, exists := round.Teams[index].GatheringPlayers[playerID]
	if !exists {
		return fmt.Errorf("player %s does not belong to team %d", playerID, index)
	}

	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/ActualPlayers/" + playerID)
	if err != nil {
		return err
	}

	err = ref.Set(player)
	if err != nil {
		round.Teams[index].ActualPlayers[playerID] = player
	}
	return err
}

// AddTeamMemberToMissing adds gathering player to missing
func (rr *RoundRepository) AddTeamMemberToMissing(round *switchers.Round, index int, playerID string) error {
	player, exists := round.Teams[index].GatheringPlayers[playerID]
	if !exists {
		return fmt.Errorf("player %s does not belong to team %d", playerID, index)
	}

	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/MissingPlayers/" + playerID)
	if err != nil {
		return err
	}

	err = ref.Set(player)
	if err != nil {
		round.Teams[index].MissingPlayers[playerID] = player
	}
	return err
}

// SetTeamState sets team state
func (rr *RoundRepository) SetTeamState(round *switchers.Round, index int, state string) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/State")
	if err != nil {
		return err
	}

	err = ref.Set(state)
	if err != nil {
		round.Teams[index].State = state
	}
	return err
}

// SetTeamMissingPlayersDeadline sets team waiting deadline after quorum
func (rr *RoundRepository) SetTeamMissingPlayersDeadline(round *switchers.Round, index int, deadline time.Time) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/MissingPlayersDeadline")
	if err != nil {
		return err
	}

	err = ref.Set(deadline)
	if err != nil {
		round.Teams[index].MissingPlayersDeadline = deadline
	}
	return err
}

// SetTeamActualDeadline sets deadline for actual task
func (rr *RoundRepository) SetTeamActualDeadline(round *switchers.Round, index int, deadline time.Time) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/ActualDeadline")
	if err != nil {
		return err
	}

	err = ref.Set(deadline)
	if err != nil {
		round.Teams[index].ActualDeadline = deadline
	}
	return err
}

// SetTeamAnswer sets team answer for actual task
func (rr *RoundRepository) SetTeamAnswer(round *switchers.Round, index int, answer *switchers.Answer) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/AnswerData")
	if err != nil {
		return err
	}

	err = ref.Set(answer)
	if err != nil {
		round.Teams[index].Answer = *answer
	}
	return err
}
