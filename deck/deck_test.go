package deck

import (
	"fmt"
	"testing"
)

func TestCreateDeck(t *testing.T) {
	c := make(chan *Deck)
	go Create(c)
	deck := <-c
	var card = deck.Cards[0]
	fmt.Println(card)
	if len(deck.Cards) != 52 {
		t.Fatalf("length is incorrect")
	}

}
