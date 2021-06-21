package dealer

import (
	"testing"
)

func TestCreate(t *testing.T) {
	c := make(chan *Dealer)
	go Create(c)
	dealer := <-c
	if len(dealer.Shoe.Cards) != dealer.Cut {
		t.Fatalf("length is incorrect")
	}


}
