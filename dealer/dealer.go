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

func splitandshuffle(cards []*card.Card) []*card.Card {
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


func Create() *Dealer {
	rand.Seed(time.Now().UnixNano())
	totaldecks := numberofdecks()
	cut := cutplacement(totaldecks)
	shoe := shoe.Create(totaldecks)
	dealer := Dealer{Id: uuid.New().String(), Cut: cut, Shoe: shoe} 

	return Shuffle(&dealer)
}


func Shuffle(dealer *Dealer) *Dealer {

	var shuffled []*card.Card

	//needs to be sequential (obviously)
	for i := 0; i < 100; i++ {
	  shuffled = splitandshuffle(dealer.Shoe.Cards)
	}
    
	//cut
	dealer.Shoe.Cuts = shuffled[dealer.Cut:]
	dealer.Shoe.Cards = shuffled[:dealer.Cut]

	fmt.Println(fmt.Sprintf("cards: %v, cut: %v", len(dealer.Shoe.Cards), len(dealer.Shoe.Cuts)))

	for _, value := range dealer.Shoe.Cards {
		fmt.Println(value.Name)
	}

    return dealer
}
