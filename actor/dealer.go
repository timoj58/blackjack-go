package actor

import (
	"github.com/google/uuid"

	"math/rand"

	"time"

	"fmt"

	"tabiiki.com/blackjack/model"

	"tabiiki.com/blackjack/util"
)

type Dealer struct {
	Id   string
	Shoe *model.Shoe
	Cut  int
}

func numberOfDecks() int {

	total := rand.Intn(10)
	if total != 0 {

		return total
	}
	return numberOfDecks()
}

func cutPlacement(decks int) int {
	total := rand.Intn(52 * decks)
	lower := (decks * 52) / 2
	upper := (decks * 52) - 10
	if total >= lower && total <= upper {
		return total
	}
	return cutPlacement(decks)
}

func (dealer *Dealer) shuffle() {

	var shuffled []*model.Card

	//needs to be sequential (obviously)
	for i := 0; i < 100; i++ {
		shuffled = util.SplitAndShuffle(dealer.Id, dealer.Shoe.Cards)
	}

	//cut
	//fmt.Println(fmt.Sprintf("cut is %v for dealer %s, shuffled length %v", dealer.Cut, dealer.Id, len(shuffled)))

	dealer.Shoe.Cuts = shuffled[dealer.Cut:]
	dealer.Shoe.Cards = shuffled[:dealer.Cut]

	fmt.Println(fmt.Sprintf("cards: %v, cut: %v for %s", len(dealer.Shoe.Cards), len(dealer.Shoe.Cuts), dealer.Id))

}

func (dealer *Dealer) reShuffle() {
	rand.Seed(time.Now().UnixNano())
	dealer.Shoe.Cards = append(dealer.Shoe.Cards, dealer.Shoe.Cuts...)
	dealer.Shoe.Cuts = dealer.Shoe.Cuts[:0]
	dealer.Cut = cutPlacement(len(dealer.Shoe.Cards) / 52)

	dealer.shuffle()
}

func CreateDealer(output chan *Dealer) {
	rand.Seed(time.Now().UnixNano())
	totaldecks := numberOfDecks()
	cut := cutPlacement(totaldecks)
	shoe := model.CreateShoe(totaldecks)
	dealer := Dealer{Id: uuid.New().String(), Cut: cut, Shoe: shoe}
	//fmt.Println(fmt.Sprintf("dealer %s, cut: %v, total cards: %v", dealer.Id, dealer.Cut, len(dealer.Shoe.Cards)))
	dealer.shuffle()
	output <- &dealer
}

func (dealer *Dealer) Hit() *model.Card {
	card := dealer.Shoe.Cards[:1][0]
	dealer.Shoe.Cards = dealer.Shoe.Cards[1:]
	dealer.Shoe.Cuts = append(dealer.Shoe.Cuts, card)
	return card
}

func Validate(cards []*model.Card) map[string]int {
	validated := make(map[string]int)
	values := util.Values(cards)

	for _, value := range values {
		if value == 21 {
			validated["Blackjack"] = value
		}
		if value < 21 {
			validated["Continue"] = value
		}
		if value > 21 {
			validated["Bust"] = value
		}
	}

	return validated
}
