package dealer

import (
	"testing"
)

func TestCreate(t *testing.T) {
	var dealer = Create()
	if len(dealer.Shoe.Cards) != dealer.Cut {
		t.Fatalf("length is incorrect")
	}


}
