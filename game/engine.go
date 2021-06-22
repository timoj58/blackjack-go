package game

import (
	"fmt"
	"tabiiki.com/blackjack/actor"
	"tabiiki.com/blackjack/model"
)

type PlayerState struct {
	Player   *actor.Player
	State    string
	Notified bool
}

type GameState struct {
	SeatingOrder []*PlayerState
	CurrentTurn  int
}

func ProcessNatural(table *Table) {
	var blackjack []string

	for _, player := range table.Players {
		if Process(player.Cards) == "Blackjack" {
			blackjack = append(blackjack, player.Id)
			broadcast(table, nil, fmt.Sprintf("player %s has blackjack", player.Id))
		}
	}

	if len(blackjack) > 0 {
		broadcast(table, nil, fmt.Sprintf("dealer: shoecard %s", table.HouseCards[1].Name))
		if Process(table.HouseCards) == "Blackjack" {
			broadcast(table, nil, fmt.Sprintf("dealer has blackjack"))
			for _, player := range blackjack {
				broadcast(table, nil, fmt.Sprintf("player %s has tied", player))
			}
		} else {
			for _, player := range blackjack {
				broadcast(table, nil, fmt.Sprintf("player %s has won", player))
			}
		}
		table.Inplay = false
	} else {
		broadcast(table, nil, "please wait for your turn to be called...")
	}
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
