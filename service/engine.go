package service

import (
	
	 "tabiiki.com/blackjack/model"
	 "tabiiki.com/blackjack/actor"
	 
)


type PlayerState struct {
	Player *actor.Player
	State string
	Notified bool
}

type GameState struct {
	SeatingOrder []*PlayerState
	CurrentTurn int
}

func Process(cards []*model.Card) string {

   	validated := actor.Validate(cards)
    if _, ok := validated["Blackjack"]; ok {
	     return "Blackjack"
	}

	if _, ok := validated["Continue"]; ok {
		return "Continue"
    }

	if _, ok := validated["Bust"]; ok {
		return "Bust"
    }

	return "Continue"

}

func Init(table *Table) {
	gameState := GameState{CurrentTurn: 0}
	for _, player := range table.Players {
        gameState.SeatingOrder = append(gameState.SeatingOrder, &PlayerState{
			Player: player, State: "Continue"})
	}

	table.GameState = &gameState
}

func NextPlayer(gameState *GameState) *actor.Player {
	return gameState.SeatingOrder[gameState.CurrentTurn].Player
}

func SetNotified(gameState *GameState, notified bool) {
	gameState.SeatingOrder[gameState.CurrentTurn].Notified = notified
}

func GetNotified(gameState *GameState) bool {
	return gameState.SeatingOrder[gameState.CurrentTurn].Notified
}




