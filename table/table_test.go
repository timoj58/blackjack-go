package table

import (
	"testing"
	"tabiiki.com/player"
	"fmt"
)



func TestCreate(t *testing.T) {
	c := make(chan *Table)
	go Create(c)
	table := <-c
	if len(table.Dealer.Shoe.Cards) != table.Dealer.Cut {
		t.Fatalf("length is incorrect")
	}

}


func TestJoin(t *testing.T) {
	c := make(chan *Table)
	go Create(c)
	table := <-c
	Join(table, player.Create(100))

	if len(table.Players) != 1 {
		t.Fatalf("players is incorrect")
	}

}

func TestStart(t *testing.T) {
	c := make(chan *Table)
	go Create(c)
	table := <-c
	Join(table, player.Create(100))
	Join(table, player.Create(100))
	Join(table, player.Create(100))
	Join(table, player.Create(100))
	Join(table, player.Create(100))

	Start(table)

	//print our the table.....
    fmt.Print(fmt.Sprintf("dealer cards: %s, ", table.HouseCards[0].Name))
	fmt.Println(fmt.Sprintf("dealer hole cards: %s", table.HouseCards[1].Name))

	for _, player := range table.Players {
		fmt.Println(fmt.Sprintf("player %s cards: %s, %s", player.Id, player.Cards[0].Name, player.Cards[1].Name))
	}

}

