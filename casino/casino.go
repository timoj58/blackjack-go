package casino


import (
	"tabiiki.com/table"

)

type Casino struct {
	Tables []*table.Table
}

func Create(tables int) *Casino {
	casino := Casino{}

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