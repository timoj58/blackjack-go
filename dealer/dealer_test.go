package dealer

import (
	"testing"
	"fmt"
)

func TestCreate(t *testing.T) {
	c := make(chan *Dealer)
	go Create(c)
	dealer := <-c
	if len(dealer.Shoe.Cards) != dealer.Cut {
		t.Fatalf("length is incorrect")
	}
}


func TestHit(t *testing.T) {
	c := make(chan *Dealer)
	go Create(c)
	dealer := <-c

	var cards = len(dealer.Shoe.Cards)
	var cuts = len(dealer.Shoe.Cuts)
	

	card := Hit(dealer)

	fmt.Println(card.Name)

	if len(dealer.Shoe.Cards) != cards - 1 {
		t.Fatalf("cards is incorrect")
	}
	if len(dealer.Shoe.Cuts) != cuts + 1 {
		t.Fatalf("cuts is incorrect")
	}

  
}	


