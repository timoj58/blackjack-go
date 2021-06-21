package card

import (
	"fmt"
)

type Card struct {
	Suit  string
	Name  string
	Value int
	Visible bool
}

func Create(suit string, name string, value int) *Card {
	card := Card{Suit: suit, Name: fmt.Sprintf("%v of %s", name, suit), Value: value, Visible: true}
	//fmt.Println(card)
	return &card
}

func Ace(suit string) *Card {
	return &Card{Suit: suit, Name: fmt.Sprintf("Ace of %s", suit), Value: 11, Visible: true}
}
