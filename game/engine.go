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

func getCard(table *Table, player *actor.Player) {
	card := table.Dealer.Hit()
	player.Cards = append(player.Cards, card)
	table.broadcast(nil, fmt.Sprintf("player %s %s", player.Id, card.Name))
}

//end of game
func Process(table *Table) {
	table.Inplay = false
}

func ProcessPlayer(table *Table, player *actor.Player) {
	switch ProcessCards(player.Cards) {
	case "Blackjack":
		table.broadcast(nil, fmt.Sprintf("player %s has blackjack", player.Id))
		Blackjack(table)
	case "Continue":
		SetPlayerState(table.GameState, "Continue")
	case "Bust":
		table.broadcast(nil, fmt.Sprintf("player %s is bust", player.Id))
		Bust(table)

	}
}

func ProcessNatural(table *Table) {
	var blackjack []string

	for _, player := range table.Players {
		if ProcessCards(player.Cards) == "Blackjack" {
			blackjack = append(blackjack, player.Id)
			table.broadcast(nil, fmt.Sprintf("player %s has blackjack", player.Id))
		}
	}

	if len(blackjack) > 0 {
		table.broadcast(nil, fmt.Sprintf("dealer hole %s", table.HouseCards[1].Name))
		if ProcessCards(table.HouseCards) == "Blackjack" {
			table.broadcast(nil, fmt.Sprintf("dealer has blackjack"))
			for _, player := range blackjack {
				table.broadcast(nil, fmt.Sprintf("player %s has tied", player))
			}
		} else {
			for _, player := range blackjack {
				table.broadcast(nil, fmt.Sprintf("player %s has won", player))
			}
		}
		table.Inplay = false
	} else {
		table.broadcast(nil, "please wait for your turn to be called...")
	}
}

func ProcessCards(cards []*model.Card) string {

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

func Start(table *Table) {
	Init(table)
	table.Inplay = true
	table.Countdown = 0
	table.HouseCards = table.HouseCards[:0]
	for _, player := range table.Players {
		player.Cards = player.Cards[:0]
	}
	//first card
	for _, player := range table.Players {
		getCard(table, player)
	}

	//dealer card
	dealerCard := table.Dealer.Hit()
	table.HouseCards = append(table.HouseCards, dealerCard)
	table.broadcast(nil, fmt.Sprintf("dealer %s", dealerCard.Name))

	//second cards
	for _, player := range table.Players {
		getCard(table, player)
	}
	//dealer hole card
	holeCard := table.Dealer.Hit()
	holeCard.Visible = false
	table.HouseCards = append(table.HouseCards, holeCard)

	ProcessNatural(table)
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

func SetPlayerState(gameState *GameState, state string) {
	gameState.SeatingOrder[gameState.CurrentTurn].State = state

}

func PlayerFinished(table *Table) {
	if table.GameState.CurrentTurn == len(table.GameState.SeatingOrder)-1 {
		//the end.  process it all.
		Process(table)
	} else {
		table.GameState.CurrentTurn++
	}
}

func Hit(table *Table, id string) {
	var player = table.Players[id]
	getCard(table, player)

	ProcessPlayer(table, player)
	SetNotified(table.GameState, false)
}

func Bust(table *Table) {
	SetPlayerState(table.GameState, "Bust")
	SetNotified(table.GameState, false)
	PlayerFinished(table)
}

func Blackjack(table *Table) {
	SetPlayerState(table.GameState, "Blackjack")
	SetNotified(table.GameState, false)
	PlayerFinished(table)
}

func Stick(table *Table, id string) {
	var player = table.Players[id]
	getCard(table, player)

	SetPlayerState(table.GameState, "Stick")
	SetNotified(table.GameState, false)
	PlayerFinished(table)
}
