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
func (table *Table) process() {
	table.supervisor.update(false)
}

func (table *Table) processPlayer(player *actor.Player) {
	switch ProcessCards(player.Cards) {
	case "Blackjack":
		table.broadcast(nil, fmt.Sprintf("player %s has blackjack", player.Id))
		table.blackjack()
	case "Continue":
		table.GameState.setPlayerState("Continue")
	case "Bust":
		table.broadcast(nil, fmt.Sprintf("player %s is bust", player.Id))
		table.bust()

	}
}

func (table *Table) processNatural() {
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
		table.supervisor.update(false)
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

func (table *Table) init() {
	gameState := GameState{CurrentTurn: 0}
	for _, player := range table.Players {
		gameState.SeatingOrder = append(gameState.SeatingOrder, &PlayerState{
			Player: player, State: "Continue"})
	}

	table.GameState = &gameState
}

func (table *Table) start() {
	table.init()
	table.Countdown = 30
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

	table.processNatural()
	table.supervisor.update(true)

}

func (gameState *GameState) nextPlayer() *actor.Player {
	return gameState.SeatingOrder[gameState.CurrentTurn].Player
}

func (gameState *GameState) setNotified(notified bool) {
	gameState.SeatingOrder[gameState.CurrentTurn].Notified = notified
}

func (gameState *GameState) getNotified() bool {
	return gameState.SeatingOrder[gameState.CurrentTurn].Notified
}

func (gameState *GameState) setPlayerState(state string) {
	gameState.SeatingOrder[gameState.CurrentTurn].State = state

}

func (table *Table) playerFinished() {
	if table.GameState.CurrentTurn == len(table.GameState.SeatingOrder)-1 {
		//the end.  process it all.
		table.process()
	} else {
		table.GameState.CurrentTurn++
	}
}

func (table *Table) hit(id string) {
	var player = table.Players[id]
	getCard(table, player)

	table.processPlayer(player)
	table.GameState.setNotified(false)
}

func (table *Table) bust() {
	table.GameState.setPlayerState("Bust")
	table.GameState.setNotified(false)
	table.playerFinished()
}

func (table *Table) blackjack() {
	table.GameState.setPlayerState("Blackjack")
	table.GameState.setNotified(false)
	table.playerFinished()
}

func (table *Table) stick(id string) {
	var player = table.Players[id]
	getCard(table, player)

	table.GameState.setPlayerState("Stick")
	table.GameState.setNotified(false)
	table.playerFinished()
}
