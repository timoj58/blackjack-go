module tabiiki.com/table

go 1.16

replace tabiiki.com/card => ../card

replace tabiiki.com/dealer => ../dealer

replace tabiiki.com/player => ../player

replace tabiiki.com/shoe => ../shoe

replace tabiiki.com/deck => ../deck

require (
	github.com/google/uuid v1.2.0
	tabiiki.com/dealer v0.0.0-00010101000000-000000000000
	tabiiki.com/player v0.0.0-00010101000000-000000000000
)
