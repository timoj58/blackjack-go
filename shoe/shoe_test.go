package shoe

import (
	"testing"
)



func TestCreate(t *testing.T) {
	var shoe = Create(10)
	if len(shoe.Cards) != (52 * 10) {
		t.Fatalf("length is incorrect")
	}

}
