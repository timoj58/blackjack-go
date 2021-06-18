package deck

import (
	"testing"
)

func TestShuffle(t *testing.T) {
	shoe := Shuffle(5)
	if len(shoe.cards) != (52 * 5) {
		t.Fatalf("length is incorrect")
	}
}

func TestCreate(t *testing.T) {
	deck := Create()
	if len(deck.cards) != 52 {
		t.Fatalf("length is incorrect")
	}

}
