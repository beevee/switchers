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
	round.FinishTime = time.Now()
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

// SaveTeam saves a single team inside a round
func (rr *RoundRepository) SaveTeam(round *switchers.Round, index int, team switchers.Team) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID + "/Teams/" + strconv.FormatInt(int64(index), 10))
	if err != nil {
		return err
	}

	return ref.Set(team)
}
