package dropsim

import (
	"log"
	"math"
	"math/rand"
)

// DropSample wraps a map of items+quantities for results sample
type DropSample struct {
	vals map[string]int
	//mu   sync.RWMutex
}

// Roll a random number in the range of the table (for concurrency test)
func (dt *DropList) roll(c chan int, x int) {
	c <- rand.Intn(x)
}

// Add drop(quantity) to sample results
func (ds *DropSample) add(key string, num int) {
	//ds.mu.Lock()
	//defer ds.mu.Unlock()
	ds.vals[key] += num
}

// TODO Concurrency testing...
func (ds *DropSample) val(key string) int {
	//ds.mu.Lock()
	//defer ds.mu.Unlock()
	return ds.vals[key]
}

// Sample uses a droptable to randomly select drops and adds the respective quantity to the output
func (dt *DropList) Sample(n int, tbl []Drop, boss string) map[string]int {
	itemCounts := DropSample{vals: make(map[string]int)}
	for _, drop := range dt.Drops {
		itemCounts.vals[drop.Name] = 0
	}
	// Roll n times * boss rolls i.e. Zulrah has 2 rolls so must roll twice per kill
	for i := 0; i < n*dt.Rolls; i++ {
		// TODO add concurrency support
		//dt.roll(c, len(tbl)-1)
		//r := <-c
		//log.Println(r)
		r := rand.Intn(len(tbl) - 1)
		itemCounts.add(tbl[r].Name, tbl[r].QuantityAvg)

	}
	// Add guaranteed drops
	for _, j := range dt.Drops {
		if j.RawRarity == 1 {
			itemCounts.vals[j.Name] += n * j.QuantityAvg
		}
	}
	// Add any special cases e.g. items dropped together
	for key, val := range SpecialCases[boss] {
		for i, val2 := range itemCounts.vals {
			if val == i {
				itemCounts.vals[key] = val2
			}
		}
	}
	return itemCounts.vals
}

// CheckTotalP Debug func for total P error
func (dt *DropList) CheckTotalP() float64 {
	total := 0.0
	for _, v := range dt.Drops {
		if v.RawRarity != 1 {
			total += v.RawRarity
		}
	}
	return total
}

// getRarest finds the item with the lowest drop rate for normalisation
func (dt *DropList) getRarest() int {
	minVal := math.Inf(1)
	for _, drop := range dt.Drops {
		if drop.RawRarity < minVal {
			minVal = drop.RawRarity
		}
	}
	return int(1/minVal) + 1
}

// MakeDropTable creates a table representing each drop as a chunk of elements, weighted by rarity
func (dt *DropList) MakeDropTable() []Drop {
	var table []Drop
	rarest := dt.getRarest()
	log.Printf("rarest = %v\n", rarest)
	for _, drop := range dt.Drops {
		if _, ok := SpecialCases[dt.BossName][drop.Name]; ok {
			// debug statement
			log.Printf("Drop %v has a special case %v\n", drop.Name, SpecialCases[dt.BossName][drop.Name])
			continue
		}
		if drop.RawRarity == 1 {
			continue
		}
		weighted := int(math.Floor(drop.RawRarity * float64(rarest)))
		var item []Drop
		for i := 0; i < weighted; i++ {
			item = append(item, drop)
		}
		table = append(table, item...)

	}
	// debug statement
	log.Printf("droptable length = %v\n", len(table))
	return table
}
