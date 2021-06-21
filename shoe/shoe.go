package shoe

import (
	"tabiiki.com/card"
	"tabiiki.com/deck"
)

type Shoe struct {
	Cards []*card.Card
	Cuts  []*card.Card
}


func Create(total int) *Shoe {

	shoe := Shoe{}
	//fix for now at 6
	c := make(chan *deck.Deck)


	for i := 0; i < total; i++ {
		go deck.Create(c)
	}

	for i := 0; i < total; i++ {
		cards := <-c
		shoe.Cards = append(shoe.Cards, cards.Cards...)
	}

	return &shoe

}
