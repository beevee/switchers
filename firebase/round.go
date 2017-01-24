package firebase

import (
	"github.com/zabawaba99/firego"

	"github.com/beevee/switchers"
)

// RoundRepository persists round information in Firebase
type RoundRepository struct {
	FirebaseURL   string
	FirebaseToken string
	firebase      *firego.Firebase
}

// Start initializes firebase connection
func (rr *RoundRepository) Start() error {
	rr.firebase = firego.New(rr.FirebaseURL, nil)
	rr.firebase.Auth(rr.FirebaseToken)

	return nil
}

// Stop does nothing
func (rr *RoundRepository) Stop() error {
	return nil
}

// CreateActiveRound creates new active round
func (rr *RoundRepository) CreateActiveRound() (*switchers.Round, error) {
	round := &switchers.Round{
		ID: "active",
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

// SaveRound saves round
func (rr *RoundRepository) SaveRound(round *switchers.Round) error {
	ref, err := rr.firebase.Ref("rounds/" + round.ID)
	if err != nil {
		return err
	}

	return ref.Set(round)
}
