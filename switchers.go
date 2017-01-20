package switchers

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(int64) (*Player, bool, error)
}

// Player is a player
type Player struct {
	ChatID int64
}
