package main

import (
	"log"
	"tabiiki.com/blackjack/game"
	"net/http"
)


func main() {
	c := game.CreateCasino(9)
	go c.Run()

	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeWs(c, w, r)
	})
    log.Fatal(http.ListenAndServe("localhost:8080", nil))
}
