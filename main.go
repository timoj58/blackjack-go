package main

import (
	"flag"
	"log"
	"net/http"
	"tabiiki.com/blackjack/game"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	c := game.CreateCasino(10)
	go c.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		game.ServeWs(c, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
