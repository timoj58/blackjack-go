package service


import (
	"net/http"
	"fmt"
	"encoding/json"
	"tabiiki.com/blackjack/actor"
	
)

type Casino struct {
	Tables map[string]*Table
	clients map[string]*Client
	broadcast chan []byte
	register chan *Client
	unregister chan *Client
}

type Message struct {
	PlayerId string  `json:"playerId"`
	Action string `json:"action"`
	Data string `json:"data"`
}


func join(casino *Casino, payload *Message) {
	//ideally would use a map for table...
	for _, t := range casino.Tables {
		if t.Id == payload.Data {
			fmt.Println(fmt.Sprintf("player %s is joining table %s", payload.PlayerId ,payload.Data))
            Join(t, casino.clients[payload.PlayerId].player)
		}
	}
}

func leave(casino *Casino, payload *Message) {
	//ideally would use a map for table...
	for _, t := range casino.Tables {
		if t.Id == payload.Data {
			fmt.Println(fmt.Sprintf("player %s is leaving table %s", payload.PlayerId ,payload.Data))
            Leave(t, casino.clients[payload.PlayerId].player)
		}
	}
	listtables(casino.clients[payload.PlayerId], casino)
}

func listtables(client *Client, casino *Casino) {
	for _, table := range casino.Tables {
        if !table.Inplay { 
		 client.send <-  []byte(fmt.Sprintf("table %s, stake: %v, %v players currently", table.Id, table.Stake, len(table.Players)))
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
		go t.Run()
		casino.Tables[t.Id] = t
	}


	return &casino
}

func (casino *Casino) Run() {
	for {
		select {
		case client := <-casino.register:
			casino.clients[client.player.Id] = client
			client.send <-  []byte(fmt.Sprintf("Welcome player %s, select a table to join", client.player.Id))
			listtables(client, casino)
		case client := <-casino.unregister:
			if _, ok := casino.clients[client.player.Id]; ok {
				delete(casino.clients, client.player.Id)
				close(client.send)
			}
		case message := <-casino.broadcast:
		    payload := Message{}
			json.Unmarshal(message, &payload)
			switch payload.Action {
			case "join":
				join(casino, &payload)
			case "leave":
				leave(casino, &payload)
			default:
				//table message received, relay it, action not relevant here.
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
		conn: conn, 
		player: actor.CreatePlayer(1000, channel),
		send: channel}
	client.casino.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}