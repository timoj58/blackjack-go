package game

import (
	"tabiiki.com/blackjack/actor"
)

type PlayerState struct {
	Player   *actor.Player
	State    string
	Notified bool
}

type GameState struct {
	SeatingOrder []*PlayerState
	CurrentTurn  int
	DealerState  *PlayerState
}

func (gameState *GameState) getHighScore() int {
	var highScore = 0

	for _, playerState := range gameState.SeatingOrder {
		if playerState.State == "Blackjack" {
			highScore = 21
		} else if playerState.State == "Stick" && highScore != 21 {
			total := actor.Validate(playerState.Player.Cards)["Continue"]
			if total > highScore {
				highScore = total
			}
		}
	}
	return highScore
}

func (gameState *GameState) nextPlayer() *actor.Player {
	return gameState.SeatingOrder[gameState.CurrentTurn].Player
}

func (gameState *GameState) setPlayerState(state string) {
	gameState.SeatingOrder[gameState.CurrentTurn].State = state

}
