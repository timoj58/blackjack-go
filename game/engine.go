package game

import (
	"fmt"
	"tabiiki.com/blackjack/actor"
	"tabiiki.com/blackjack/model"
)

func (table *Table) getCard(player *actor.Player) {
	card := table.Dealer.Hit()
	player.Cards = append(player.Cards, card)
	table.broadcast(nil, fmt.Sprintf("player %s %s", player.Id, card.Name))
}

func (table *Table) getDealerCard() {
	dealerCard := table.Dealer.Hit()
	table.HouseCards = append(table.HouseCards, dealerCard)
	table.broadcast(nil, fmt.Sprintf("dealer %s", dealerCard.Name))

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
				table.broadcast(nil, fmt.Sprintf("player %s wins", playerState.Player.Id))
			} else {
				playerState.Player.Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("player %s ties", playerState.Player.Id))
			}
		} else if playerState.State == "Stick" && highScore != 21 {
			total := actor.Validate(playerState.Player.Cards)["Continue"]
			if total == highScore && total > dealerTotal {
				playerState.Player.Funds += table.Stake * 2
				table.broadcast(nil, fmt.Sprintf("player %s wins", playerState.Player.Id))
			} else if total == highScore && dealerTotal == total {
				playerState.Player.Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("player %s ties", playerState.Player.Id))
			}
		}
	}
}

//end of game
func (table *Table) process() {

	table.broadcast(nil, "Dealer to play...")
	//dealer needs to show
	table.broadcast(nil, fmt.Sprintf("dealer hole %s", table.HouseCards[1].Name))

	var highScore = table.GameState.getHighScore()
	table.processDealer(highScore)
	var dealerTotal = table.getDealerTotal()

	for _, playerState := range table.GameState.SeatingOrder {
		table.broadcast(nil, fmt.Sprintf("player %s state %s", playerState.Player.Id, playerState.State))
	}
	table.broadcast(nil, fmt.Sprintf("dealer state %s", table.GameState.DealerState.State))

	table.processWinners(highScore, dealerTotal)

	if dealerTotal > highScore {
		table.broadcast(nil, "dealer wins")
	}

	table.supervisor.update(false)
}

func (table *Table) processPlayer(player *actor.Player) {
	switch ProcessCards(player.Cards) {
	case "Blackjack":
		table.broadcast(nil, fmt.Sprintf("player %s has blackjack", player.Id))
		table.blackjack()
	case "Continue":
		table.GameState.setPlayerState("Continue")
		table.broadcast(table.GameState.nextPlayer(), "Its your turn!")
	case "Bust":
		table.broadcast(nil, fmt.Sprintf("player %s is bust", player.Id))
		table.bust()

	}
}

func (table *Table) processDealer(highScore int) {
	switch ProcessCards(table.HouseCards) {
	case "Blackjack":
		table.broadcast(nil, "dealer has blackjack")
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
		table.broadcast(nil, "dealer is is bust")
		table.GameState.DealerState.State = "Bust"
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
				table.Players[player].Funds += table.Stake
				table.broadcast(nil, fmt.Sprintf("player %s has tied", player))
			}
		} else {
			for _, player := range blackjack {
				table.Players[player].Funds += table.Stake * 2
				table.broadcast(nil, fmt.Sprintf("player %s has won", player))
			}
		}
		table.supervisor.update(false)
	} else {
		table.broadcast(nil, "please wait for your turn to be called...")
		table.broadcast(table.GameState.nextPlayer(), "Its your turn!")
		table.supervisor.update(true)
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
		table.broadcast(table.GameState.nextPlayer(), "Its your turn!")
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
	table.GameState.setPlayerState("Stick")
	table.playerFinished()
}
