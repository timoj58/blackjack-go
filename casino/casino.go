package casino


import (
	"tabiiki.com/table"
	"net/http"
	"fmt"
	"encoding/json"
	"tabiiki.com/player"

)

type Casino struct {
	Tables []*table.Table
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
            table.Join(t, casino.clients[payload.PlayerId].player)
		}
	}
}

func Create(tables int) *Casino {
	casino := Casino{
		broadcast:  make(chan []byte),
		register:   make(chan *Client),
		unregister: make(chan *Client),
		clients:    make(map[string]*Client),
	}

	c := make(chan *table.Table)
    
	for i := 0; i < tables; i++ {
		go table.Create(c)
	}

	for i := 0; i < tables; i++ {
		table := <-c
		casino.Tables = append(casino.Tables, table)
	}


	return &casino
}

func (casino *Casino) Run() {
	for {
		select {
		case client := <-casino.register:
			casino.clients[client.player.Id] = client
			client.send <-  []byte(fmt.Sprintf("Welcome player %s, select a table to join", client.player.Id))
			for _, table := range casino.Tables {
				client.send <-  []byte(fmt.Sprintf("table %s, %v players currently", table.Id, len(table.Players)))
			}
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
				//leave table..todo
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
	client := &Client{
		casino: casino, 
		conn: conn, 
		player: player.Create(1000),
		send: make(chan []byte, 256)}
	client.casino.register <- client

	// Allow collection of memory referenced by the caller by doing all work in
	// new goroutines.
	go client.writePump()
	go client.readPump()
}