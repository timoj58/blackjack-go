package table

import (
	"github.com/google/uuid"

	"tabiiki.com/dealer"

	"tabiiki.com/card"

	"tabiiki.com/player"

)


type Table struct {
   Id string
   Dealer *dealer.Dealer
   Players []*player.Player
   HouseCards []*card.Card
}

func Create(output chan *Table) {
	c := make(chan *dealer.Dealer)
	go dealer.Create(c)
	table := Table{Id: uuid.New().String(), Dealer : <- c}
	output <- &table
}

func Join(table *Table, player *player.Player) {
    table.Players = append(table.Players, player)
}

func Leave(table *Table, player *player.Player) {

}

func Start(table *Table) {
	table.HouseCards = table.HouseCards[:0]
	for _, player := range table.Players {
		player.Cards = player.Cards[:0]
	}
	//first card
	for _, player := range table.Players {
		player.Cards = append(player.Cards, dealer.Hit(table.Dealer))
	}

	//dealer card 
	table.HouseCards = append(table.HouseCards, dealer.Hit(table.Dealer))
	//second cards
	for _, player := range table.Players {
		player.Cards = append(player.Cards, dealer.Hit(table.Dealer))		
	}
	//dealer hole card
	card := dealer.Hit(table.Dealer)
    card.Visible = false
	table.HouseCards = append(table.HouseCards, card)

}