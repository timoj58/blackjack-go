package actor

import (
	"github.com/google/uuid"
   "tabiiki.com/blackjack/model"

)

type Player struct {
   Id string
   Funds int
   Cards []*model.Card
}

func CreatePlayer(funds int) *Player {
   player := Player{Id: uuid.New().String(), Funds: funds}
   return &player
}