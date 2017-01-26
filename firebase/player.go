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

// SetState sets player state
func (pr *PlayerRepository) SetState(player *switchers.Player, state string) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/State")
	if err != nil {
		return err
	}

	err = ref.Set(state)
	if err == nil {
		player.State = state
	}
	return err
}

// SetName sets player name
func (pr *PlayerRepository) SetName(player *switchers.Player, name string) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/Name")
	if err != nil {
		return err
	}

	err = ref.Set(name)
	if err == nil {
		player.Name = name
	}
	return err
}

// SetPaused sets player paused
func (pr *PlayerRepository) SetPaused(player *switchers.Player, paused bool) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/Paused")
	if err != nil {
		return err
	}

	err = ref.Set(paused)
	if err == nil {
		player.Paused = paused
	}
	return err
}

// SetTrump sets player Trump
func (pr *PlayerRepository) SetTrump(player *switchers.Player, trump bool) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/Trump")
	if err != nil {
		return err
	}

	err = ref.Set(trump)
	if err == nil {
		player.Trump = trump
	}
	return err
}

// SetModeratingTeamIndex sets player Trump
func (pr *PlayerRepository) SetModeratingTeamIndex(player *switchers.Player, index int) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/ModeratingTeamIndex")
	if err != nil {
		return err
	}

	err = ref.Set(index)
	if err == nil {
		player.ModeratingTeamIndex = index
	}
	return err
}

// IncreaseScore increases player score
func (pr *PlayerRepository) IncreaseScore(player *switchers.Player) error {
	ref, err := pr.firebase.Ref("players/" + player.ID + "/Score")
	if err != nil {
		return err
	}

	err = ref.Set(player.Score + 1)
	if err == nil {
		player.Score++
	}
	return err
}
