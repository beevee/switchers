package gameprocessor

import (
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

const (
	playerStateNew        = ""
	playerStateAskName    = "askname"
	playerStateIdle       = "idle"
	playerStateGathering  = "gathering"
	playerStatePlaying    = "playing"
	playerStateModeration = "moderation"

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
		player.Trump = true
	}

	if player.Trump {
		gp.executeTrumpCommand(command, player)
	} else {
		gp.executePlayerCommand(command, player)
	}

	if err := gp.PlayerRepository.SavePlayer(player); err != nil {
		gp.Logger.Log("msg", "error saving player profile after executing command", "playerid", player.ID)
		gp.Bot.SendMessage(player.ID, "Что-то пошло не так, возможно надо попробовать еще раз.")
	}
}
