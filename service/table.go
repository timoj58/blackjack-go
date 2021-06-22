package service



import (
	
	"fmt"
	 "tabiiki.com/blackjack/actor"
	 "tabiiki.com/blackjack/model"
	 "github.com/google/uuid"
	 "time"
	 "math/rand"

)


type Table struct {
   Id string
   Dealer *actor.Dealer
   Players map[string]*actor.Player
   HouseCards []*model.Card
   Countdown int
   Inplay bool
   Stake int
   GameState *GameState
}

func tablestake() int {
	stakes := []int{10, 25, 50, 75}
	index := rand.Intn(4)

	return stakes[index]
}


func broadcast(table *Table, sender *actor.Player, message string) {
	for _, player := range table.Players {
        if sender == nil || sender.Id == player.Id {
		 player.Send <- []byte(message)
		}
	}
}

func processNatural(table *Table) {
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
 	}else{
		for _, player := range blackjack {
			broadcast(table, nil, fmt.Sprintf("player %s has won", player))
 		}
 	}
	table.Inplay = false
  }else{
	broadcast(table, nil, "please wait for your turn to be called...")
  }
}


func CreateTable(output chan *Table) {
	rand.Seed(time.Now().UnixNano())

	c := make(chan *actor.Dealer)
	go actor.CreateDealer(c)
	table := Table{
		Id: uuid.New().String(), 
		Dealer : <-c, 
		Players: make(map[string]*actor.Player),
		Stake: tablestake(),
		Inplay: false}
	output <- &table
}

func Join(table *Table, player *actor.Player) {
	fmt.Println(fmt.Sprintf("player %s has joined", player.Id))
    table.Players[player.Id] = player
	if len(table.Players) == 1 {
		table.Countdown = 10
	}
	//player.Send <- []byte("hello there waiting for new players")
}

func Leave(table *Table, player *actor.Player) {
	fmt.Println(fmt.Sprintf("player %s has left", player.Id))
    delete(table.Players,player.Id)
}

func Start(table *Table) {
	Init(table)
	table.Inplay = true
	table.HouseCards = table.HouseCards[:0]
	for _, player := range table.Players {
		player.Cards = player.Cards[:0]
	}
	//first card
	for _, player := range table.Players {
		card := actor.Hit(table.Dealer)
		player.Cards = append(player.Cards, card)
		broadcast(table, nil, fmt.Sprintf("player %s: card %s", player.Id, card.Name))
	}

	//dealer card 
	dealerCard := actor.Hit(table.Dealer)
	table.HouseCards = append(table.HouseCards, dealerCard)
	broadcast(table, nil, fmt.Sprintf("dealer: card %s", dealerCard.Name))

	//second cards
	for _, player := range table.Players {
		secondCard := actor.Hit(table.Dealer)
		player.Cards = append(player.Cards, secondCard)		
		broadcast(table, nil, fmt.Sprintf("player %s: card %s", player.Id, secondCard.Name))
	}
	//dealer hole card
	holeCard := actor.Hit(table.Dealer)
    holeCard.Visible = false
	table.HouseCards = append(table.HouseCards, holeCard)

	processNatural(table)

}

func (table *Table) Run() {
	
	for {
        if !table.Inplay && len(table.Players) > 0 {

            if table.Countdown == 0 {
				broadcast(table, nil, "game starting...")
                Start(table)
			}else{
		    	//countdown till start
			    time.Sleep(time.Second)

				broadcast(table, nil, fmt.Sprintf("%v seconds till game starts", table.Countdown))

			    table.Countdown -= 1
	    	}

		}

		if table.Inplay {
            
			time.Sleep(2 * time.Second)
            if !GetNotified(table.GameState) {
			 broadcast(table, NextPlayer(table.GameState), "Its your turn")
			 SetNotified(table.GameState, true)
			}
			
		}
       
	}
}
	
