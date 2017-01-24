package gameprocessor

import (
	"time"

	"github.com/beevee/switchers"
)

// PopulateRound populates a new round according to the rules of game
func (gp *GameProcessor) PopulateRound(round *switchers.Round) error {
	round.StartTime = time.Now()

	return nil
}
