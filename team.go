package switchers

import "time"

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
	Answer                 Answer
}

// Answer is a team answer with original message details
type Answer struct {
	Text      string
	OwnerID   string
	MessageID string
}

// IsEmpty detects if Answer contains usable information
func (a Answer) IsEmpty() bool {
	return a.Text == "" && a.MessageID == ""
}
