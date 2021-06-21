package player

import (
	"github.com/google/uuid"
	
	"tabiiki.com/card"
)

type Player struct {
   Id string
   Funds int
   Cards []*card.Card
}

func Create(funds int) *Player {
   player := Player{Id: uuid.New().String(), Funds: funds}
   return &player
}