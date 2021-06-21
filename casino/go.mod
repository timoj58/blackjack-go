module tabiiki.com/casino

go 1.16

replace tabiiki.com/table => ../table

replace tabiiki.com/dealer => ../dealer

replace tabiiki.com/player => ../player

replace tabiiki.com/shoe => ../shoe

replace tabiiki.com/deck => ../deck

replace tabiiki.com/card => ../card

require (
	github.com/gorilla/websocket v1.4.2
	tabiiki.com/player v0.0.0-00010101000000-000000000000
	tabiiki.com/table v0.0.0-00010101000000-000000000000
)
