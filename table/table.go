package table

import (
	"github.com/google/uuid"

	"tabiiki.com/dealer"

	"tabiiki.com/player"

)


type Table struct {
   Id string
   Dealer *dealer.Dealer
   Players []*player.Player
}

func Create(output chan *Table) {
	c := make(chan *dealer.Dealer)
	go dealer.Create(c)
	table := Table{Id: uuid.New().String(), Dealer : <- c}
	output <- &table
}

func Join(table *Table, player *player.Player) {

}

func Leave(table *Table, player *player.Player) {

}