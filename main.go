package main

import (
	"flag"
	"log"
	"net/http"
	"tabiiki.com/casino"
)

var addr = flag.String("addr", ":8080", "http service address")


func main() {
	flag.Parse()
	c := casino.Create(10)
	go c.Run()
	http.HandleFunc("/ws", func(w http.ResponseWriter, r *http.Request) {
		casino.ServeWs(c, w, r)
	})
	err := http.ListenAndServe(*addr, nil)
	if err != nil {
		log.Fatal("ListenAndServe: ", err)
	}
}
