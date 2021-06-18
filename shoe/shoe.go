package shoe

import (
	"tabiiki.com/card"
	"tabiiki.com/deck"
)

type Shoe struct {
	Cards []*card.Card
	Cuts  []*card.Card
}

func singledeck(c chan *deck.Deck) {
	go deck.Create(c)
}

func fourdecks(c chan *deck.Deck) {

	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)

}

func eightdecks(c chan *deck.Deck) {

	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)

	//output <- <-c, <-c, <-c, <-c, <-c, <-c, <-c, <-c

}

func defaultdecks(output chan *deck.Deck) {
	c := make(chan *deck.Deck)

	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)
	go deck.Create(c)

	//output <- <-c, <-c, <-c, <-c, <-c
}

func Create(total int) *Shoe {

	shoe := Shoe{}
	//fix for now at 6
	c := make(chan *deck.Deck)

	switch total {
	case 1:
		singledeck(c)
		x := <-c
		shoe.Cards = append(shoe.Cards, x.Cards...)
	case 4:
		fourdecks(c)
		x, y, z, q := <-c, <-c, <-c, <-c
		shoe.Cards = append(shoe.Cards, x.Cards...)
		shoe.Cards = append(shoe.Cards, y.Cards...)
		shoe.Cards = append(shoe.Cards, z.Cards...)
		shoe.Cards = append(shoe.Cards, q.Cards...)
	case 8:
		eightdecks(c)
	default:
		defaultdecks(c)
	}

	return &shoe
}
