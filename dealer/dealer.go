package dealer

import (
	"github.com/google/uuid"

	"tabiiki.com/shoe"

)

type Dealer struct {
	Id   string
	Shoe shoe.Shoe
	Cut  int //default cut rating mode... allow players to influence dealer on shuffle.
}

func Create() *Dealer {
	dealer := Dealer{Id: uuid.New().String(), Cut: 1} //cut should be random value.  id should be UUID
	return &dealer
}

func Start() *shoe.Shoe {
	return nil
}

func Shuffle(shoe *shoe.Shoe) *shoe.Shoe {
    return nil
}
