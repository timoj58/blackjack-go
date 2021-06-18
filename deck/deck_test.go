package deck

import (
	"fmt"
	"testing"
)

func TestCreateShoe(t *testing.T) {
	shoe := CreateShoe()
	if len(shoe.cards) != (52 * 5) {
		t.Fatalf("length is incorrect")
	}
}

func TestCreateDeck(t *testing.T) {
	c := make(chan *Deck)
	go CreateDeck(c)
	deck := <-c
	var card = deck.cards[0]
	fmt.Println(card)
	if len(deck.cards) != 52 {
		t.Fatalf("length is incorrect")
	}

}
