package model

type Shoe struct {
	Cards []*Card
	Cuts  []*Card
}

func CreateShoe(total int) *Shoe {

	shoe := Shoe{}
	c := make(chan *Deck)

	for i := 0; i < total; i++ {
		go CreateDeck(c)
	}

	for i := 0; i < total; i++ {
		cards := <-c
		shoe.Cards = append(shoe.Cards, cards.Cards...)
	}

	return &shoe

}
