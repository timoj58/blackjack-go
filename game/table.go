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
	table.broadcast(nil, fmt.Sprintf("player %s has joined table", player.Id))
	table.Players[player.Id] = player
}

func (table *Table) leave(player *actor.Player) {
	delete(table.Players, player.Id)
	table.broadcast(player, fmt.Sprintf("player %s has left table", player.Id))
}

func (table *Table) event(message *Message) {
	switch message.Action {
	case "hit":
		table.hit(message.PlayerId)
	case "stick":
		table.stick(message.PlayerId)
	default:
		//ignore other cases for now. (split etc)
	}
}

func (table *Table) run() {

	for {
		inplay := <-table.supervisor.c

		if !inplay && len(table.Players) > 0 {

			if table.Countdown == 0 {
				table.broadcast(nil, "game starting...")
				table.start()
			} else {
				time.Sleep(time.Second)

				table.broadcast(nil, fmt.Sprintf("%v seconds till game starts...", table.Countdown))
				table.Countdown -= 1
			}

		}

	}

}
