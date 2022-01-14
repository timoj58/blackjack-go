module tabiiki.com/blackjack

go 1.16

require (
	github.com/google/uuid v1.2.0
	github.com/gorilla/websocket v1.4.2
	github.com/rs/cors v1.8.2 // indirect
)

replace tabiiki.com/blackjack/actor => ../actor

replace tabiiki.com/blackjack/model => ../model

replace tabiiki.com/blackjack/util => ../util

replace tabiiki.com/blackjack/game => ./game
