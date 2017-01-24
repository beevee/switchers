package gameprocessor

import (
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

const (
	playerStateNew     = ""
	playerStateAskName = "askname"
	playerStateIdle    = "idle"

	teamMinSize = 6

	teamStateGathering  = "gathering"
	teamStatePlaying    = "playing"
	teamStateModeration = "moderation"
	teamStateWon        = "won"
	teamStateLost       = "lost"

	commandNewRound = "/newround"
	commandResign   = "/resign"

	commandSetName = "/setname"
	commandPause   = "/pause"
	commandResume  = "/resume"
)

// GameProcessor contains all in-game logic
type GameProcessor struct {
	TrumpCode        string
	PlayerRepository switchers.PlayerRepository
	RoundRepository  switchers.RoundRepository
	TaskRepository   switchers.TaskRepository
	Bot              switchers.Bot
	Logger           switchers.Logger
	tomb             tomb.Tomb
}

// Start initializes loops that make game go round
func (gp *GameProcessor) Start() error {
	gp.tomb.Go(gp.roundDeactivator)
	gp.tomb.Go(gp.deadlineEnforcer)

	return nil
}

// Stop gracefully stops loops
func (gp *GameProcessor) Stop() error {
	gp.tomb.Kill(nil)
	return gp.tomb.Wait()
}
