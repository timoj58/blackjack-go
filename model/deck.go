package model

import (
	"fmt"
)

type Deck struct {
	Cards []*Card
}

func numbercards(suit string, values [9]int, c chan []*Card) {
	var cards []*Card
	for _, value := range values {
		cards = append(cards, CreateCard(suit, fmt.Sprintf("%v", value), value))
	}
	c <- cards

}

func facecards(suit string, values [3]string, c chan []*Card) {
	var cards []*Card
	for _, value := range values {
		cards = append(cards, CreateCard(suit, value, 10))
	}
	c <- cards

}

func suitcards(suit string, output chan []*Card) {
	c := make(chan []*Card)
	var cards []*Card
	values := [9]int{2, 3, 4, 5, 6, 7, 8, 9, 10}
	faces := [3]string{"King", "Queen", "Jack"}

	go numbercards(suit, values, c)
	go facecards(suit, faces, c)
	x, y := <-c, <-c
	cards = append(cards, x...)
	cards = append(cards, y...)
	cards = append(cards, Ace(suit))

	output <- cards

}

func CreateDeck(output chan *Deck) {
	deck := Deck{}
	c := make(chan []*Card)

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
