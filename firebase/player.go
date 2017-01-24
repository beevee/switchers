package firebase

import "github.com/beevee/switchers"

// PlayerRepository persists player information in Firebase
type PlayerRepository struct {
	Repository
}

// GetOrCreatePlayer retrieves player profile, creating it if necessary
func (pr *PlayerRepository) GetOrCreatePlayer(ID string) (*switchers.Player, bool, error) {
	ref, err := pr.firebase.Ref("players/" + ID)
	if err != nil {
		return nil, false, err
	}

	player := &switchers.Player{}
	created := false
	if err = ref.Value(player); err != nil {
		return nil, false, err
	}
	if player.ID != ID {
		player.ID = ID
		if err = ref.Set(player); err != nil {
			return nil, false, err
		}
		created = true
	}

	return player, created, nil
}

// GetAllPlayers retrieves all players
func (pr *PlayerRepository) GetAllPlayers() (map[string]*switchers.Player, error) {
	ref, err := pr.firebase.Ref("players")
	if err != nil {
		return nil, err
	}

	var players map[string]*switchers.Player
	if err = ref.Value(&players); err != nil {
		return nil, err
	}
	return players, nil
}

// GetAllTrumps retrieves all Trumps
func (pr *PlayerRepository) GetAllTrumps() (map[string]*switchers.Player, error) {
	ref, err := pr.firebase.Ref("players")
	if err != nil {
		return nil, err
	}

	var players map[string]*switchers.Player
	if err = ref.OrderBy("Trump").EqualTo("true").Value(&players); err != nil {
		return nil, err
	}
	return players, nil
}

// SavePlayer saves player profile
func (pr *PlayerRepository) SavePlayer(player *switchers.Player) error {
	ref, err := pr.firebase.Ref("players/" + player.ID)
	if err != nil {
		return err
	}

	return ref.Set(player)
}
