package switchers

import "time"

// Round is a round
type Round struct {
	ID        string
	StartTime time.Time
	Teams     []*Team
}
