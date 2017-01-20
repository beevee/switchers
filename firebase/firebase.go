package firebase

import (
	"strconv"

	"github.com/zabawaba99/firego"

	"github.com/beevee/switchers"
)

// PlayerRepository persists player information in Firebase
type PlayerRepository struct {
	FirebaseURL   string
	FirebaseToken string
	firebase      *firego.Firebase
}

// Start initializes firebase connection
func (pr *PlayerRepository) Start() error {
	pr.firebase = firego.New(pr.FirebaseURL, nil)
	pr.firebase.Auth(pr.FirebaseToken)

	return nil
}

// Stop does nothing
func (pr *PlayerRepository) Stop() error {
	return nil
}

// GetOrCreatePlayer retrieves player profile, creating it if necessary
func (pr *PlayerRepository) GetOrCreatePlayer(chatID int64) (*switchers.Player, bool, error) {
	ref, err := pr.firebase.Ref("users/" + strconv.FormatInt(chatID, 10))
	if err != nil {
		return nil, false, err
	}

	player := &switchers.Player{}
	created := false
	if err = ref.Value(player); err != nil {
		return nil, false, err
	}
	if player.ChatID != chatID {
		player.ChatID = chatID
		if err = ref.Set(player); err != nil {
			return nil, false, err
		}
		created = true
	}

	return player, created, nil
}
