package actor

import (
	"fmt"
	"testing"
)

func TestShuffle(t *testing.T) {
	c := make(chan *Dealer)
	go CreateDealer(c)

	dealer := <-c

	if len(dealer.Shoe.Cards) != dealer.Cut {
		t.Fatalf("players is incorrect")
	}

	for _, card := range dealer.Shoe.Cards {
		fmt.Println(card.Name)
	}

}
