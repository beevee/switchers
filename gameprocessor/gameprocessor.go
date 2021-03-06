package gameprocessor

import (
	"gopkg.in/tomb.v2"

	"github.com/beevee/switchers"
)

// GameProcessor contains all in-game logic
type GameProcessor struct {
	TrumpCode             string
	TeamQuorum            int
	TeamMinSize           int
	PlayerRepository      switchers.PlayerRepository
	RoundRepository       switchers.RoundRepository
	TaskRepository        switchers.TaskRepository
	Bot                   switchers.Bot
	Logger                switchers.Logger
	playerCommandChannels map[string]chan command
	tomb                  tomb.Tomb
}

type command struct {
	id   string
	text string
}

// Start initializes loops that make game go round
func (gp *GameProcessor) Start() error {
	gp.playerCommandChannels = make(map[string]chan command)

	gp.tomb.Go(gp.gameProgressor)

	return nil
}

// Stop gracefully stops loops
func (gp *GameProcessor) Stop() error {
	gp.tomb.Kill(nil)
	return gp.tomb.Wait()
}

// ExecuteCommand takes text command from a player and schedules it for execution in a separate channel for exch user
func (gp *GameProcessor) ExecuteCommand(commandID string, commandText string, playerID string) {
	_, exists := gp.playerCommandChannels[playerID]

	if !exists {
		gp.playerCommandChannels[playerID] = make(chan command, 1000)
		gp.tomb.Go(func() error {
			gp.playerCommandExecutor(playerID, gp.playerCommandChannels[playerID])
			return nil
		})
		gp.Logger.Log("msg", "created new command executor goroutine for player", "playerid", playerID)
	}

	cmd := command{
		id:   commandID,
		text: commandText,
	}
	gp.playerCommandChannels[playerID] <- cmd
}

func (gp *GameProcessor) playerCommandExecutor(playerID string, playerCommands <-chan command) {
	for {
		select {
		case cmd := <-playerCommands:
			player, created, err := gp.PlayerRepository.GetOrCreatePlayer(playerID)
			if err != nil {
				gp.Logger.Log("msg", "error retrieving player profile while executing command", "error", err)
			}
			if created {
				gp.Logger.Log("msg", "created player profile", "playerid", player.ID)
			}

			if cmd.text == gp.TrumpCode {
				if err := gp.PlayerRepository.SetTrump(player, true); err != nil {
					gp.Bot.SendMessage(player.ID, responseSomethingWrong)
					gp.Logger.Log("msg", "failed to make player Trump", "error", err)
					return
				}
			}

			if player.Trump {
				gp.executeTrumpCommand(cmd, player)
			} else {
				gp.executePlayerCommand(cmd, player)
			}

		case <-gp.tomb.Dying():
			gp.Logger.Log("msg", "aborted command executor goroutine for player", "playerid", playerID)
			return
		}
	}
}
