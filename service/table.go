package service



import (
	
	"fmt"
	 "tabiiki.com/blackjack/actor"
	 "tabiiki.com/blackjack/model"
	 "github.com/google/uuid"
	

)


type Table struct {
   Id string
   Dealer *actor.Dealer
   Players map[string]*actor.Player
   HouseCards []*model.Card
   Countdown int
   Inplay bool
}

func CreateTable(output chan *Table) {
	c := make(chan *actor.Dealer)
	go actor.CreateDealer(c)
	table := Table{
		Id: uuid.New().String(), 
		Dealer : <-c, 
		Players: make(map[string]*actor.Player),
		Inplay: false}
	output <- &table
}

func Join(table *Table, player *actor.Player) {
	fmt.Println(fmt.Sprintf("player %s has joined", player.Id))
    table.Players[player.Id] = player
}

func Leave(table *Table, player *actor.Player) {
	fmt.Println(fmt.Sprintf("player %s has left", player.Id))
    delete(table.Players,player.Id)
}

func Start(table *Table) {
	table.Inplay = true
	table.HouseCards = table.HouseCards[:0]
	for _, player := range table.Players {
		player.Cards = player.Cards[:0]
	}
	//first card
	for _, player := range table.Players {
		player.Cards = append(player.Cards, actor.Hit(table.Dealer))
	}

	//dealer card 
	table.HouseCards = append(table.HouseCards, actor.Hit(table.Dealer))
	//second cards
	for _, player := range table.Players {
		player.Cards = append(player.Cards, actor.Hit(table.Dealer))		
	}
	//dealer hole card
	card := actor.Hit(table.Dealer)
    card.Visible = false
	table.HouseCards = append(table.HouseCards, card)

}

func (table *Table) Run() {
	
	for {
		/*select {
		} */
	}
}
	
