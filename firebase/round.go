package firebase

import (
	"strconv"
	"time"

	"github.com/beevee/switchers"
	"github.com/satori/go.uuid"
)

// RoundRepository persists round information in Firebase
type RoundRepository struct {
	Repository
}

// CreateActiveRound creates new active round
func (rr *RoundRepository) CreateActiveRound() (*switchers.Round, error) {
	round := &switchers.Round{
		ID:        "active",
		StartTime: time.Now(),
	}

	ref, err := rr.firebase.Ref("rounds/active")
	if err != nil {
		return nil, err
	}

	if err = ref.Set(round); err != nil {
		return nil, err
	}

	return round, nil
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
	if err := rr.SaveRound(round); err != nil {
		return err
	}

	ref, err := rr.firebase.Ref("rounds/active")
	if err != nil {
		return err
	}

	return ref.Remove()
}

// SaveRound saves round
func (rr *RoundRepository) SaveRound(round *switchers.Round) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID)
	if err != nil {
		return err
	}

	return ref.Set(round)
}

// SetPlayerGathered sets flag that a single player in team has gathered
func (rr *RoundRepository) SetPlayerGathered(round *switchers.Round, index int, playerID string) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/PlayerIDs/" + playerID)
	if err != nil {
		return err
	}

	err = ref.Set(true)
	if err != nil {
		round.Teams[index].PlayerIDs[playerID] = true
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
func (rr *RoundRepository) SetTeamAnswer(round *switchers.Round, index int, answer string) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10) + "/Answer")
	if err != nil {
		return err
	}

	err = ref.Set(answer)
	if err != nil {
		round.Teams[index].Answer = answer
	}
	return err
}
