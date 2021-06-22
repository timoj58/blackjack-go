package main

import (
	"flag"
	"log"
	"net/http"
	"tabiiki.com/blackjack/service"
)

var addr = flag.String("addr", ":8080", "http service address")

func main() {
	flag.Parse()
	c := service.CreateCasino(10)
	go c.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		service.ServeWs(c, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
