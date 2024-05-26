package main

import (
	"fmt"
	"math"
	"math/rand"
)

func (dt *DropTable) Sample(n int) {
	itemCounts := make(map[string]int)
	for _, v := range dt.Drops {
		itemCounts[v.Name] = 0
	}
	table := dt.MakeTable()
	for i := 0; i < n; i++ {
		r := rand.Intn(len(table))
		itemCounts[table[r].Name] += table[r].QuantityAvg
	}
	fmt.Printf("From %d kills:\n", n)
	for k, v := range itemCounts {
		if v > 0 {
			fmt.Printf("You received %dx %s\n", v, k)
		}

	}
}

func (dt *DropTable) CheckTotalP() float64 {
	total := 0.0
	for _, v := range dt.Drops {
		if v.RawRarity != 1 {
			total += v.RawRarity
		}
	}
	return total
}

func (dt *DropTable) getRarest() int {
	minVal := math.Inf(1)
	for _, drop := range dt.Drops {
		if drop.RawRarity < minVal {
			minVal = drop.RawRarity
		}
	}
	return int(1 / minVal)
}

func (dt *DropTable) MakeTable() []Drop {
	var table []Drop
	for _, drop := range dt.Drops {
		if drop.RawRarity == 1 {
			continue
		}
		weighted := int(math.Floor(drop.RawRarity * float64(dt.getRarest())))
		var item []Drop
		for i := 0; i < weighted; i++ {
			item = append(item, drop)
		}
		table = append(table, item...)
		//fmt.Printf("item: %s, weighted: %d\n", drop.Name, weighted)

	}
	return table
}
