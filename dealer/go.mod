module tabiiki.com/dealer

go 1.16

replace tabiiki.com/shoe => ../shoe

require (
	github.com/google/uuid v1.2.0
	tabiiki.com/card v0.0.0-00010101000000-000000000000
	tabiiki.com/shoe v0.0.0-00010101000000-000000000000
)

replace tabiiki.com/card => ../card

replace tabiiki.com/deck => ../deck
