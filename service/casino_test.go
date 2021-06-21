package service

import (
	"testing"
)



func TestCreate(t *testing.T) {
	var casino = CreateCasino(10)
	if len(casino.Tables) != 10 {
		t.Fatalf("length is incorrect")
	}

}
