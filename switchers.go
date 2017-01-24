package switchers

import "time"

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(ID string) (*Player, bool, error)
	SavePlayer(player *Player) error
}

// Player is a player
type Player struct {
	ID     string
	Trump  bool
	State  string
	Name   string
	Paused bool
	Score  int
}

// RoundRepository persists round information
type RoundRepository interface {
	CreateActiveRound() (*Round, error)
	GetActiveRound() (*Round, error)
	SaveRound(round *Round) error
}

// Round is a round
type Round struct {
	ID        string
	StartTime time.Time
}

// GameProcessor contains all in-game logic
type GameProcessor interface {
	ExecuteCommand(string, *Player) string
}
