package game

import (
	"fmt"
	"github.com/google/uuid"
	"math/rand"
	"tabiiki.com/blackjack/actor"
	"tabiiki.com/blackjack/model"
	"time"
)

type Table struct {
	Id         string
	Dealer     *actor.Dealer
	Players    map[string]*actor.Player
	HouseCards []*model.Card
	Countdown  int
	supervisor *TableSupervisor
	Stake      int
	GameState  *GameState
}

func tableStake() int {
	stakes := []int{10, 25, 50, 75}
	index := rand.Intn(4)

	return stakes[index]
}


func (table *Table) broadcast(sender *actor.Player, message string) {
	for _, player := range table.Players {
		if sender == nil || sender.Id == player.Id {
			defer func() {
				if err := recover(); err != nil {
					fmt.Println("panic occurred:", err)
				}
			}()  
			player.Send <- []byte(message)
		}
	}
}

func CreateTable(output chan *Table) {
	rand.Seed(time.Now().UnixNano())

	c := make(chan *actor.Dealer)
	go actor.CreateDealer(c)
	table := Table{
		Id:         uuid.New().String(),
		Dealer:     <-c,
		Players:    make(map[string]*actor.Player),
		Stake:      tableStake(),
		Countdown:  15,
		supervisor: CreateTableSupervisor(make(chan bool))}

	go table.supervisor.run()

	output <- &table
}

func (table *Table) join(player *actor.Player) {
	if len(table.Players) < 7 {
		table.broadcast(nil, fmt.Sprintf("player %s has joined table", player.Id))
		table.Players[player.Id] = player
	} else {
		table.broadcast(player, "table is full, try another")
	}
}

func (table *Table) leave(player *actor.Player) {
	delete(table.Players, player.Id)
	table.broadcast(nil, fmt.Sprintf("{\"type\": \"players\",\"data\": \"player %s has left table\", \"id\": \"%s\"}", player.Id, player.Id))
}

func (table *Table) event(message *Message) {
	if message.PlayerId == table.GameState.currentPlayer().Id {
		switch message.Action {
		case "hit":
			table.hit(message.PlayerId)
		case "stick":
			table.stick(message.PlayerId)
		default:
			//ignore other cases for now. (split etc)
		}
	} else {
		table.broadcast(table.Players[message.PlayerId], "its not your turn")
	}
}

func (table *Table) checkFunds() {
	for _, player := range table.Players {
		if player.Funds < table.Stake {
			table.broadcast(player, "you have insufficient funds for this table, you have been kicked out!!!")
			table.leave(player)
		}
	}
}

func (table *Table) run() {

	for {
		select {
		case inplay := <-table.supervisor.c:

		if !inplay && len(table.Players) > 0 {

			//make sure the players have the funds...else kick them out.
			table.checkFunds()
			if table.Countdown == 0 {
				table.broadcast(nil, "{\"type\": \"message\", \"data\": \"game starting...\"}")
				table.start()
			} else {
				table.broadcast(nil, fmt.Sprintf("{\"type\": \"message\", \"data\": \"%v seconds till game starts...\"}", table.Countdown))
				table.Countdown -= 1
				time.Sleep(time.Second)
			}

		} else if !inplay && len(table.Players) == 0 {
			//reset incase we have kicked some players out, or they left
			table.Countdown = 15
		}
	}

	}

}
