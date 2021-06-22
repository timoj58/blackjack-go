package model

import (
	"testing"
	"fmt"
)



func TestNode(t *testing.T) {

	var cards []*Card

	cards  = append(cards, CreateCard("hearts", "3 of hearts", 3))
	cards  = append(cards, CreateCard("clubs", "ace of clubs", 11))
	cards  = append(cards, CreateCard("spades", "ace of spades", 11))
	cards  = append(cards, CreateCard("hearts", "2 of hearts", 2))

	var hand = CreateHand(cards)
	

	for _, h := range hand {
		fmt.Println(fmt.Sprintf("values %v", h))
	}

}
