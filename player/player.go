package player

import (
	"github.com/google/uuid"
)

type Player struct {
   Id string
}

func Create() *Player {
   player := Player{Id: uuid.New().String()}
   return &player
}