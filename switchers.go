package switchers

import "time"

// PlayerRepository persists player information
type PlayerRepository interface {
	GetOrCreatePlayer(ID string) (*Player, bool, error)
	GetAllPlayers() (map[string]*Player, error)
	GetAllTrumps() (map[string]*Player, error)
	SetState(player *Player, state string) error
	SetPaused(player *Player, paused bool) error
	SetTrump(player *Player, trump bool) error
	SetName(player *Player, name string) error
	IncreaseScore(player *Player) error
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
	GetActiveRound() (*Round, error)
	DeactivateRound(*Round) error
	SaveActiveRound(round *Round) error
	AddTeamMemberToActual(round *Round, index int, playerID string) error
	AddTeamMemberToMissing(round *Round, index int, playerID string) error
	SetTeamState(round *Round, index int, state string) error
	SetTeamMissingPlayersDeadline(round *Round, index int, deadline time.Time) error
	SetTeamActualDeadline(round *Round, index int, deadline time.Time) error
	SetTeamAnswer(round *Round, index int, answer string) error
}

// Round is a round
type Round struct {
	ID        string
	StartTime time.Time
	Teams     []*Team
}

// Team is a team
type Team struct {
	State                  string
	GatheringPlayers       map[string]Player
	ActualPlayers          map[string]Player
	MissingPlayers         map[string]Player
	GatheringTask          GatheringTask
	GatheringDeadline      time.Time
	MissingPlayersDeadline time.Time
	ActualTask             ActualTask
	ActualDeadline         time.Time
	Answer                 string
}

// TaskRepository persists task information
type TaskRepository interface {
	GetAllGatheringTasks() ([]GatheringTask, error)
	GetAllActualTasks() ([]ActualTask, error)
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
	CorrectAnswers   []string
}

// Bot maintains communication with players
type Bot interface {
	SendMessage(ID string, message string)
}

// GameProcessor contains all in-game logic
type GameProcessor interface {
	ExecuteCommand(command string, playerID string)
}
