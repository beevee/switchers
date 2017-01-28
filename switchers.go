package switchers

import "time"

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(ID string) (*Player, bool, error)
	GetAllPlayers() (map[string]*Player, error)
	GetAllTrumps() (map[string]*Player, error)
	GetTop(count int) ([]*Player, error)
	SetState(player *Player, state string) error
	SetPaused(player *Player, paused bool) error
	SetTrump(player *Player, trump bool) error
	SetName(player *Player, name string) error
	SetModeratingTeamIndex(player *Player, index int) error
	IncreaseScore(player *Player) error
}

// RoundRepository persists round information
type RoundRepository interface {
	GetActiveRound() (*Round, error)
	DeactivateRound(*Round) error
	SaveActiveRound(round *Round) error
	AddTeamMemberToActual(round *Round, index int, playerID string) error
	AddTeamMemberToMissing(round *Round, index int, playerID string) error
	SetTeamState(round *Round, index int, state string) error
	SetTeamMissingPlayersDeadline(round *Round, index int, deadline time.Time) error
	SetTeamActualDeadline(round *Round, index int, deadline time.Time) error
	SetTeamAnswer(round *Round, index int, answer *Answer) error
}

// TaskRepository persists task information
type TaskRepository interface {
	GetAllGatheringTasks() ([]GatheringTask, error)
	GetAllActualTasks() ([]ActualTask, error)
}

// Bot maintains communication with players
type Bot interface {
	SendMessage(ID string, message string)
	ForwardMessage(ID string, messageText string, messageID string, messageOwnerID string)
}

// GameProcessor contains all in-game logic
type GameProcessor interface {
	ExecuteCommand(commandID string, commandText string, playerID string)
}
