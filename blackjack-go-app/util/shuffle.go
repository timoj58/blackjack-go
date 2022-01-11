package util

import (
	"math/rand"

	"tabiiki.com/blackjack/model"
)

func shufflesplit(cards []*model.Card, c chan []*model.Card) {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
	c <- cards
}

func SplitAndShuffle(dealer string, cards []*model.Card, output chan []*model.Card) {

	var shuffled []*model.Card
	c := make(chan []*model.Card)

	go shufflesplit(cards[:len(cards)/2], c)
	go shufflesplit(cards[len(cards)/2:], c)

	x, y := <-c, <-c

	shuffled = append(shuffled, y...)
	shuffled = append(shuffled, x...)

	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	output <- shuffled

}
