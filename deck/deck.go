package deck

import (
	"fmt"
)

type Card struct {
	suit  string
	name  string
	value int
}

type Deck struct {
	cards []Card
}

type Shoe struct {
	cards []Card
	cut   int
}

func Create() *Deck {
	values := [9]int{1, 2, 3, 4, 5, 6, 7, 8, 9}
	faces := [3]string{"King", "Queen", "Jack"}
	suits := [4]string{"Diamonds", "Hearts", "Spades", "Clubs"}
	deck := Deck{}

	for _, suit := range suits {
		for _, value := range values {
			deck.cards = append(deck.cards, Card{suit: suit, name: fmt.Sprintf("%v of %s", value, suit), value: value})
		}

		for _, face := range faces {
			deck.cards = append(deck.cards, Card{suit: suit, name: fmt.Sprintf("%s of %s", face, suit), value: 10})
		}
		deck.cards = append(deck.cards, Card{suit: suit, name: fmt.Sprintf("Ace of %s", suit), value: 11})
	}

	return &deck
}

func Shuffle(total int) *Shoe {
	var decks []*Deck
	shoe := Shoe{cut: 26 * total} //improve this.  ie dealer based.

	for i := 0; i < total; i++ {
		decks = append(decks, Create())
	}

	//ignore shuffle and just add all the cards.
	for _, deck := range decks {
		shoe.cards = append(shoe.cards, deck.cards...)
	}

	return &shoe
}
