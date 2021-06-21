package dealer

import (
	"github.com/google/uuid"

	"tabiiki.com/shoe"

	"tabiiki.com/card"

	"math/rand"

	"time"

	"fmt"

)


type Dealer struct {
	Id   string
	Shoe *shoe.Shoe
	Cut  int 
}

func numberofdecks() int {

   total := rand.Intn(10)
   if(total != 0) {

	   return total
   }
   return numberofdecks()
}

func cutplacement(decks int) int {
	total := rand.Intn(52 * decks)
	lower := (decks * 52) /2
	upper := (decks * 52) - 10
	if(total >= lower && total <= upper) {
		return total
	}
	return cutplacement(decks)
}

func splitandshuffle(dealer string, cards []*card.Card) []*card.Card {
	
	//fmt.Println(fmt.Sprintf("dealer %s is shuffling", dealer))
	var shuffled []*card.Card
    c := make(chan []*card.Card)

	go shufflesplit(cards[:len(cards)/2], c)
	go shufflesplit(cards[len(cards)/2:], c)

	x, y := <-c, <-c
		
	shuffled = append(shuffled, y...)
    shuffled = append(shuffled, x...)

	rand.Shuffle(len(shuffled), func(i, j int) { shuffled[i], shuffled[j] = shuffled[j], shuffled[i] })

	return shuffled

}

func shufflesplit(cards []*card.Card, c chan []*card.Card) {
	rand.Shuffle(len(cards), func(i, j int) { cards[i], cards[j] = cards[j], cards[i] })
     c <- cards
}


func shuffle(dealer *Dealer) {

	var shuffled []*card.Card

	//needs to be sequential (obviously)
	for i := 0; i < 100; i++ {
	  shuffled = splitandshuffle(dealer.Id, dealer.Shoe.Cards)
	}
    
	//cut
	//fmt.Println(fmt.Sprintf("cut is %v for dealer %s, shuffled length %v", dealer.Cut, dealer.Id, len(shuffled)))

	dealer.Shoe.Cuts = shuffled[dealer.Cut:]
	dealer.Shoe.Cards = shuffled[:dealer.Cut]

	fmt.Println(fmt.Sprintf("cards: %v, cut: %v for %s", len(dealer.Shoe.Cards), len(dealer.Shoe.Cuts), dealer.Id))

}

func Reshuffle(dealer *Dealer) {
	rand.Seed(time.Now().UnixNano())
	dealer.Shoe.Cards = append(dealer.Shoe.Cards, dealer.Shoe.Cuts...)
	dealer.Shoe.Cuts = dealer.Shoe.Cuts[:0]
	dealer.Cut = cutplacement(len(dealer.Shoe.Cards)/52)

	shuffle(dealer)
}

func Create(output chan *Dealer) {
	rand.Seed(time.Now().UnixNano())
	totaldecks := numberofdecks()
	cut := cutplacement(totaldecks)
	shoe := shoe.Create(totaldecks)
	dealer := Dealer{Id: uuid.New().String(), Cut: cut, Shoe: shoe} 
	//fmt.Println(fmt.Sprintf("dealer %s, cut: %v, total cards: %v", dealer.Id, dealer.Cut, len(dealer.Shoe.Cards)))
    shuffle(&dealer)
	output <- &dealer
}


func Hit(dealer *Dealer) *card.Card {
   card := dealer.Shoe.Cards[:1][0]
   dealer.Shoe.Cards = dealer.Shoe.Cards[1:]
   dealer.Shoe.Cuts = append(dealer.Shoe.Cuts, card)	
   return card
}

func Check(cards []*card.Card) int {
   total := 0

   for _, card := range cards {
	 total += card.Value
   }

   return total
} 
