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
	Inplay     bool //this needs to be a channel
	Stake      int
	GameState  *GameState
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

func CreateTable(output chan *Table) {
	rand.Seed(time.Now().UnixNano())

	c := make(chan *actor.Dealer)
	go actor.CreateDealer(c)
	table := Table{
		Id:      uuid.New().String(),
		Dealer:  <-c,
		Players: make(map[string]*actor.Player),
		Stake:   tablestake(),
		Inplay:  false}
	output <- &table
}

func Join(table *Table, player *actor.Player) {
	broadcast(table, nil, fmt.Sprintf("player %s has joined", player.Id))
	table.Players[player.Id] = player
	if len(table.Players) == 1 {
		table.Countdown = 10
	}
}

func Leave(table *Table, player *actor.Player) {
	delete(table.Players, player.Id)
	broadcast(table, player, fmt.Sprintf("player %s has left", player.Id))
}

func Event(table *Table, message *Message) {
	switch message.Action {
	case "hit":
		Hit(table, message.PlayerId)
	case "stick":
		Stick(table, message.PlayerId)
	default:
		//ignore other cases for now. (split etc)
	}
}

func (table *Table) Run() {

	for {
		if !table.Inplay && len(table.Players) > 0 {

			if table.Countdown == 0 {
				broadcast(table, nil, "game starting...")
				Start(table)
			} else {
				//countdown till start
				time.Sleep(time.Second)

				broadcast(table, nil, fmt.Sprintf("%v seconds till game starts...", table.Countdown))

				table.Countdown -= 1
			}

		}

		//need to put a lock on this.....
		if table.Inplay {

			time.Sleep(2 * time.Second)
			if !GetNotified(table.GameState) {
				broadcast(table, NextPlayer(table.GameState), "Its your turn!")
				SetNotified(table.GameState, true)
			}

		}

	}
}
