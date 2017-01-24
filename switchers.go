package switchers

import "time"

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(ID string) (*Player, bool, error)
	GetAllPlayers() (map[string]*Player, error)
	GetAllTrumps() (map[string]*Player, error)
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
	DeactivateRound(*Round) error
	SaveRound(round *Round) error
	SaveTeam(round *Round, index int, team Team) error
}

// Round is a round
type Round struct {
	ID         string
	StartTime  time.Time
	FinishTime time.Time
	Teams      []Team
}

// Team is a team
type Team struct {
	State             string
	PlayerIDs         map[string]bool
	GatheringTask     GatheringTask
	GatheringDeadline time.Time
	ActualTask        ActualTask
	ActualDeadline    time.Time
	Answer            string
}

// TaskRepository persists task information
type TaskRepository interface {
	GetAllGatheringTasks() ([]GatheringTask, error)
}

// GatheringTask is a task to gather team
type GatheringTask struct {
	Text             string
	TimeLimitMinutes int
}

// ActualTask is a task for team that gathered
type ActualTask struct {
	Text             string
	TimeLimitMinutes int
	CorrectAnswer    string
}

// Bot maintains communication with players
type Bot interface {
	SendMessage(ID string, message string)
}

// GameProcessor contains all in-game logic
type GameProcessor interface {
	ExecuteCommand(command string, playerID string)
}
