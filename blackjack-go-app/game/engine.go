package game

import (
	"fmt"
	"time"
	"tabiiki.com/blackjack/actor"
	"tabiiki.com/blackjack/model"
)

func (table *Table) getCard(player *actor.Player) {
	card := table.Dealer.Hit()
	player.Cards = append(player.Cards, card)
	table.broadcast(nil, fmt.Sprintf("{\"type\": \"card\", \"data\": \"%s\", \"id\": \"%s\", \"dealer\": false}", card.Name, player.Id))
	time.Sleep(time.Second)
}

func (table *Table) getDealerCard() {
	dealerCard := table.Dealer.Hit()
	table.HouseCards = append(table.HouseCards, dealerCard)
	table.broadcast(nil, fmt.Sprintf("{\"type\": \"card\", \"data\": \"%s\", \"id\": \"%s\", \"dealer\": true}", dealerCard.Name, table.Dealer.Id))
	time.Sleep(time.Second)
}

func (table *Table) getDealerTotal() int {
	if table.GameState.DealerState.State == "Blackjack" {
		return 21
	} else if table.GameState.DealerState.State == "Stick" {
		return actor.Validate(table.HouseCards)["Continue"]
	}
	return 0
}

func (table *Table) processWinners(highScore int, dealerTotal int) {
	for _, playerState := range table.GameState.SeatingOrder {
		if playerState.State == "Blackjack" {
			if dealerTotal != 21 {
				playerState.Player.Funds += table.Stake * 2
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s wins\"}", playerState.Player.Id))
				time.Sleep(time.Second)
			} else {
				playerState.Player.Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s ties\"}", playerState.Player.Id))
				time.Sleep(time.Second)
			}
		} else if playerState.State == "Stick" && highScore != 21 {
			total := actor.Validate(playerState.Player.Cards)["Continue"]
			if total == highScore && total > dealerTotal {
				playerState.Player.Funds += table.Stake * 2
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s wins\"}", playerState.Player.Id))
				time.Sleep(time.Second)
			} else if total == highScore && dealerTotal == total {
				playerState.Player.Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s ties\"}", playerState.Player.Id))
				time.Sleep(time.Second)
			}
		}
	}
}

//end of game
func (table *Table) process() {

	table.broadcast(nil, "{\"type\": \"game\", \"data\": \"Dealer to play...\"}")
	time.Sleep(time.Second)
	//dealer needs to show
	table.broadcast(nil, fmt.Sprintf("{\"type\": \"card\", \"data\": \"%s\", \"id\": \"%s\", \"dealer\": true}", table.HouseCards[1].Name, table.Dealer.Id))
	time.Sleep(time.Second)

	var highScore = table.GameState.getHighScore()
	table.processDealer(highScore)
	var dealerTotal = table.getDealerTotal()

	for _, playerState := range table.GameState.SeatingOrder {
		table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"player %s state %s\"}", playerState.Player.Id, playerState.State))
		time.Sleep(time.Second)
	}
	table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"dealer state %s\"}", table.GameState.DealerState.State))
	time.Sleep(time.Second)

	table.processWinners(highScore, dealerTotal)
	time.Sleep(time.Second)

	if dealerTotal > highScore {
		table.broadcast(nil, "{\"type\": \"result\", \"data\": \"dealer wins\"}")
		time.Sleep(time.Second)
	}

	table.supervisor.update(false)
	//TODO need to broadcast to everyone...not table.
	table.Casino.globalBroadcast(fmt.Sprintf("{\"type\": \"table-status\", \"status\": true, \"id\": \"%s\"}", table.Id))
}

func (table *Table) processPlayer(player *actor.Player) {
	switch ProcessCards(player.Cards) {
	case "Blackjack":
		table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"player %s has blackjack\"}", player.Id))
		time.Sleep(time.Second)
		table.blackjack()
		time.Sleep(time.Second)
	case "Continue":
		table.GameState.setPlayerState("Continue")
		table.broadcast(table.GameState.nextPlayer(), "{\"type\": \"game\", \"data\": \"Its your turn!\"}")
		time.Sleep(time.Second)
	case "Bust":
		table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"player %s is bust\"}", player.Id))
		time.Sleep(time.Second)
		table.bust()
		time.Sleep(time.Second)

	}
}

func (table *Table) processDealer(highScore int) {
	switch ProcessCards(table.HouseCards) {
	case "Blackjack":
		table.broadcast(nil, "{\"type\": \"game\", \"data\": \"dealer has blackjack\"}")
		time.Sleep(time.Second)
		table.GameState.DealerState.State = "Blackjack"
	case "Continue":
		total := actor.Validate(table.HouseCards)["Continue"]
		if total < highScore {
			table.getDealerCard()
			table.processDealer(highScore)
		} else {
			//dealer will stick.
			table.GameState.DealerState.State = "Stick"
		}
	case "Bust":
		table.broadcast(nil, "{\"type\": \"game\", \"data\": \"dealer is is bust\"}")
		time.Sleep(time.Second)
		table.GameState.DealerState.State = "Bust"
	}
}

func (table *Table) processNatural() {
	var blackjack []string

	for _, player := range table.Players {
		if ProcessCards(player.Cards) == "Blackjack" {
			blackjack = append(blackjack, player.Id)
			table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"player %s has blackjack\"}", player.Id))
			time.Sleep(time.Second)
		}
	}

	if len(blackjack) > 0 {
		table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"dealer hole %s\"}", table.HouseCards[1].Name))
		time.Sleep(time.Second)
		if ProcessCards(table.HouseCards) == "Blackjack" {
			table.broadcast(nil, fmt.Sprintf("{\"type\": \"game\", \"data\": \"dealer has blackjack\"}"))
			time.Sleep(time.Second)
			for _, player := range blackjack {
				table.Players[player].Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s has tied\"}", player))
				time.Sleep(time.Second)
			}
		} else {
			for _, player := range blackjack {
				table.Players[player].Funds += table.Stake * 2
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"result\", \"data\": \"player %s has won\"}", player))
				time.Sleep(time.Second)
			}
		}
		table.supervisor.update(false)
		//TODO need to broadcast to everyone...not table.
		table.Casino.globalBroadcast(fmt.Sprintf("{\"type\": \"table-status\", \"status\": true, \"id\": \"%s\"}", table.Id))
		} else {
		table.broadcast(nil, "{\"type\": \"game\", \"data\": \"please wait for your turn to be called...\"}")
		time.Sleep(time.Second)
		table.broadcast(table.GameState.nextPlayer(), "{\"type\": \"game\", \"data\": \"Its your turn!\"}")
		time.Sleep(time.Second)
		table.supervisor.update(true)
		//TODO need to broadcast to everyone...not table.
		table.Casino.globalBroadcast(fmt.Sprintf("{\"type\": \"table-status\", \"status\": false, \"id\": \"%s\"}", table.Id))
	
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

	gameState.DealerState = &PlayerState{}
	table.GameState = &gameState
}

func (table *Table) start() {
	table.init()
	table.Countdown = 15
	table.HouseCards = table.HouseCards[:0]
	for _, player := range table.Players {
		player.Cards = player.Cards[:0]
		player.Funds -= table.Stake
	}
	//first card
	for _, player := range table.Players {
		table.getCard(player)
	}

	//dealer card
	table.getDealerCard()

	//second cards
	for _, player := range table.Players {
		table.getCard(player)
	}
	//dealer hole card
	holeCard := table.Dealer.Hit()
	holeCard.Visible = false
	table.HouseCards = append(table.HouseCards, holeCard)

	table.processNatural()
}

func (table *Table) playerFinished() {
	if table.GameState.CurrentTurn == len(table.GameState.SeatingOrder)-1 {
		//the end.  process it all.
		table.process()
	} else {
		table.GameState.CurrentTurn++
		table.broadcast(table.GameState.nextPlayer(), "{\"type\": \"game\", \"data\": \"Its your turn!\"}")
		time.Sleep(time.Second)
	}
}

func (table *Table) hit(id string) {
	var player = table.Players[id]
	table.getCard(player)

	table.processPlayer(player)
}

func (table *Table) bust() {
	table.GameState.setPlayerState("Bust")
	table.playerFinished()
}

func (table *Table) blackjack() {
	table.GameState.setPlayerState("Blackjack")
	table.playerFinished()
}

func (table *Table) stick(id string) {
	var player = table.Players[id]
	table.GameState.setPlayerState("Stick")
	table.playerFinished()
	table.broadcast(player, "{\"type\": \"game\", \"data\": \"Please wait\"}")
	
}
