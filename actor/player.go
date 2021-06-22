package actor

import (
	"github.com/google/uuid"
	"tabiiki.com/blackjack/model"
)

type Player struct {
	Id    string
	Funds int
	Cards []*model.Card
	Send  chan []byte
}

func CreatePlayer(funds int, send chan []byte) *Player {
	player := Player{Id: uuid.New().String(), Funds: funds, Send: send}
	return &player
}
