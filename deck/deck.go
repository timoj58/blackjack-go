package deck

import (
	"fmt"

	"tabiiki.com/card"
)

type Deck struct {
	Cards []*card.Card
}

func numbercards(suit string, values [8]int, c chan []*card.Card) {
	var cards []*card.Card
	for _, value := range values {
		cards = append(cards, card.Create(suit, fmt.Sprintf("%v", value), value))
	}
	c <- cards

}

func facecards(suit string, values [3]string, c chan []*card.Card) {
	var cards []*card.Card
	for _, value := range values {
		cards = append(cards, card.Create(suit, value, 10))
	}
	c <- cards

}

func suitcards(suit string, output chan []*card.Card) {
	c := make(chan []*card.Card)
	var cards []*card.Card
	values := [8]int{2, 3, 4, 5, 6, 7, 8, 9}
	faces := [3]string{"King", "Queen", "Jack"}

	go numbercards(suit, values, c)
	go facecards(suit, faces, c)
	x, y := <-c, <-c
	cards = append(cards, x...)
	cards = append(cards, y...)
	cards = append(cards, card.Ace(suit))

	output <- cards

}

func Create(output chan *Deck) {
	deck := Deck{}
	c := make(chan []*card.Card)

	go suitcards("Hearts", c)
	go suitcards("Diamonds", c)
	go suitcards("Clubs", c)
	go suitcards("Spades", c)

	spades, clubs, diamonds, hearts := <-c, <-c, <-c, <-c

	deck.Cards = append(deck.Cards, spades...)
	deck.Cards = append(deck.Cards, clubs...)
	deck.Cards = append(deck.Cards, diamonds...)
	deck.Cards = append(deck.Cards, hearts...)

	output <- &deck
}
