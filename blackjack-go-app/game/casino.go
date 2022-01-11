package game

import (
	"encoding/json"
	"fmt"
	"net/http"
	"tabiiki.com/blackjack/actor"
)

type Casino struct {
	Tables     map[string]*Table
	clients    map[string]*Client
	broadcast  chan []byte
	register   chan *Client
	unregister chan *Client
}

type Message struct {
	PlayerId string `json:"playerId"`
	Action   string `json:"action"`
	Data     string `json:"data"`
}

func (casino *Casino) funds(payload *Message) {
	var client = casino.clients[payload.PlayerId]
	client.send <- []byte(fmt.Sprintf("your funds are %v", client.player.Funds))
}

func (casino *Casino) event(payload *Message) {
	var table = casino.Tables[payload.Data]
	if <-table.supervisor.c {
		table.event(payload)
	}
}

func (casino *Casino) join(payload *Message) {
	var table = casino.Tables[payload.Data]

	if !<-table.supervisor.c {
		table.join(casino.clients[payload.PlayerId].player)
	} else {
		casino.clients[payload.PlayerId].send <- []byte("table is currently in session")
	}
}

func (casino *Casino) leave(payload *Message) {
	var table = casino.Tables[payload.Data]

	if !<-table.supervisor.c {
		table.leave(casino.clients[payload.PlayerId].player)
		casino.listTables(casino.clients[payload.PlayerId])
	} else {
		casino.clients[payload.PlayerId].send <- []byte("table is currently in session")
	}
}

func (casino *Casino) listTables(client *Client) {
	for _, table := range casino.Tables {
		if !<-table.supervisor.c {
			client.send <- []byte(fmt.Sprintf("table %s, cut: %v stake: %v, %v players", table.Id, table.Dealer.Cut, table.Stake, len(table.Players)))
		}
	}
}

func CreateCasino(tables int) *Casino {
	casino := Casino{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
		Tables:     make(map[string]*Table),
	}

	c := make(chan *Table)

	for i := 0; i < tables; i++ {
		go CreateTable(c)
	}

	for i := 0; i < tables; i++ {
		t := <-c
		go t.run()
		casino.Tables[t.Id] = t
	}

	return &casino
}

func (casino *Casino) Run() {
	for {
		select {
		case client := <-casino.register:
			casino.clients[client.player.Id] = client
			client.send <- []byte(fmt.Sprintf("Welcome player %s, select a table to join", client.player.Id))
			casino.listTables(client)
		case client := <-casino.unregister:
			if _, ok := casino.clients[client.player.Id]; ok {
				delete(casino.clients, client.player.Id)
				close(client.send)
			}
		case message := <-casino.broadcast:
			payload := Message{}
			json.Unmarshal(message, &payload)
			switch payload.Action {
			case "list":
				casino.listTables(casino.clients[payload.PlayerId])
			case "join":
				casino.join(&payload)
			case "leave":
				casino.leave(&payload)
			case "funds":
				casino.funds(&payload)
			default:
				casino.event(&payload)
			}
		}
	}
}

// serveWs handles websocket requests from the peer.
func ServeWs(casino *Casino, w http.ResponseWriter, r *http.Request) {
	conn, err := upgrader.Upgrade(w, r, nil)
	if err != nil {
		return
	}
	channel := make(chan []byte, 256)
	client := &Client{
		casino: casino,
		conn:   conn,
		player: actor.CreatePlayer(1000, channel),
		send:   channel}
	client.casino.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}