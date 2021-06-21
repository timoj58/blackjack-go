package table

import (
	"testing"
)



func TestCreate(t *testing.T) {
	c := make(chan *Table)
	go Create(c)
	table := <-c
	if len(table.Dealer.Shoe.Cards) != table.Dealer.Cut {
		t.Fatalf("length is incorrect")
	}

}
