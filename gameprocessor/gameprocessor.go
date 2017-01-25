package gameprocessor

import (
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

const (
	playerStateNew     = ""
	playerStateAskName = "askname"
	playerStateIdle    = "idle"
	playerStateInGame  = "ingame"

	teamMinSize    = 6
	teamQuorumSize = 4

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

const (
	somethingWrongResponse = "Что-то пошло не так. Попробуй еще раз."
)

// ExecuteCommand takes text command from a player, updates internal state and returns response
func (gp *GameProcessor) ExecuteCommand(command string, playerID string) {
	player, created, err := gp.PlayerRepository.GetOrCreatePlayer(playerID)
	if err != nil {
		gp.Logger.Log("msg", "error retrieving player profile while executing command", "error", err)
	}
	if created {
		gp.Logger.Log("msg", "created player profile", "playerid", player.ID)
	}

	if command == gp.TrumpCode {
		if err := gp.PlayerRepository.SetTrump(player, true); err != nil {
			gp.Bot.SendMessage(player.ID, somethingWrongResponse)
			gp.Logger.Log("msg", "failed to make player Trump", "error", err)
			return
		}
	}

	if player.Trump {
		gp.executeTrumpCommand(command, player)
	} else {
		gp.executePlayerCommand(command, player)
	}
}
