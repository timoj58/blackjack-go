package deck

import (
	"fmt"
)

type Card struct {
	suit  string
	name  string
	value int
}

type Deck struct {
	cards []*Card
}

type Shoe struct {
	cards []*Card
	cut   int
}

func createcard(suit string, name string, value int) *Card {
	return &Card{suit: suit, name: fmt.Sprintf("%v of %s", name, suit), value: value}
}

func numbercards(suit string, values [9]int, c chan []*Card) {
	var cards []*Card
	for _, value := range values {
		cards = append(cards, createcard(suit, fmt.Sprintf("%v", value), value))
	}
	c <- cards

}

func facecards(suit string, values [3]string, c chan []*Card) {
	var cards []*Card
	for _, value := range values {
		cards = append(cards, createcard(suit, value, 10))
	}
	c <- cards

}

func suitcards(suit string, output chan []*Card) {
	c := make(chan []*Card)
	var cards []*Card
	values := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	faces := [3]string{"King", "Queen", "Jack"}

	go numbercards(suit, values, c)
	go facecards(suit, faces, c)
	x, y := <-c, <-c
	cards = append(cards, x...)
	cards = append(cards, y...)
	cards = append(cards, &Card{suit: suit, name: fmt.Sprintf("Ace of %s", suit), value: 11})

	output <- cards

}

func CreateDeck(output chan *Deck) {
	deck := Deck{}
	c := make(chan []*Card)

	go suitcards("Hearts", c)
	go suitcards("Diamonds", c)
	go suitcards("Aces", c)
	go suitcards("Spades", c)

	spades, aces, diamonds, hearts := <-c, <-c, <-c, <-c

	deck.cards = append(deck.cards, spades...)
	deck.cards = append(deck.cards, aces...)
	deck.cards = append(deck.cards, diamonds...)
	deck.cards = append(deck.cards, hearts...)

	output <- &deck
}

func CreateShoe() *Shoe {

	shoe := Shoe{cut: 26 * 5} //make ir random
	//fix for now at 6
	c := make(chan *Deck)

	go CreateDeck(c)
	go CreateDeck(c)
	go CreateDeck(c)
	go CreateDeck(c)
	go CreateDeck(c)

	x, y, z, j, k := <-c, <-c, <-c, <-c, <-c
	shoe.cards = append(shoe.cards, x.cards...)
	shoe.cards = append(shoe.cards, y.cards...)
	shoe.cards = append(shoe.cards, z.cards...)
	shoe.cards = append(shoe.cards, j.cards...)
	shoe.cards = append(shoe.cards, k.cards...)

	return &shoe
}

func Shuffle(shoe *Shoe) *Shoe {
	return shoe
}
