package shoe

import (
	"testing"
)

func TestCreateShoe(t *testing.T) {
	var shoe = Create(1)
	if len(shoe.Cards) != 52 {
		t.Fatalf("length is incorrect")
	}

}

func TestCreateShoe4(t *testing.T) {
	var shoe = Create(4)
	if len(shoe.Cards) != (52 * 4) {
		t.Fatalf("length is incorrect")
	}

}

func TestCreateShoe8(t *testing.T) {
	var shoe = Create(8)
	if len(shoe.Cards) != (52 * 8) {
		t.Fatalf("length is incorrect")
	}

}
