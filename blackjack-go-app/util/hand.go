package util

import (
	"tabiiki.com/blackjack/model"
)

func Values(cards []*model.Card) []int {

	var values []int
	values = append(values, 0)

	for _, card := range cards {
		for i := range values {
			if card.Value == 11 {
				values = append(values, values[i]+1)
			}
			values[i] = values[i] + card.Value
		}
	}

	return values

}
