package casino

import (
	"testing"
)



func TestCreate(t *testing.T) {
	var casino = Create(10)
	if len(casino.Tables) != 10 {
		t.Fatalf("length is incorrect")
	}

}
