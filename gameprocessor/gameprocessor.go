package gameprocessor

import "github.com/beevee/switchers"

// GameProcessor contains all in-game logic
type GameProcessor struct {
	PlayerRepository switchers.PlayerRepository
	RoundRepository  switchers.RoundRepository
}
