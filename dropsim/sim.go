package dropsim

import (
	"math"
	"math/rand"
	"sync"
)

type DropSample struct {
	vals map[string]int
	mu   sync.Mutex
}

func (dt *DropTable) roll(c chan int, x int) {
	c <- rand.Intn(x)
}

func (ds *DropSample) add(key string, num int) {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	ds.vals[key] += num
}

func (ds *DropSample) val(key string) int {
	ds.mu.Lock()
	defer ds.mu.Unlock()
	return ds.vals[key]
}

func (dt *DropTable) Sample(n int) map[string]int {
	itemCounts := DropSample{vals: make(map[string]int)}
	c := make(chan int)
	for _, v := range dt.Drops {
		itemCounts.vals[v.Name] = 0
	}
	table := dt.MakeTable()

	for i := 0; i < n*dt.Rolls; i++ {
		go dt.roll(c, len(table))
		go func() {
			r := <-c
			itemCounts.add(table[r].Name, table[r].QuantityAvg)
		}()

	}
	for _, j := range dt.Drops {
		if j.RawRarity == 1 {
			itemCounts.vals[j.Name] += n * j.QuantityAvg
		}
	}
	return itemCounts.vals
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

	}
	return table
}
