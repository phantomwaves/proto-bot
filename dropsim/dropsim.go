package dropsim

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"strconv"
	"strings"
)

var BannedItems = "brimstone key" +
	"frozen key piece (saradomin)" +
	"frozen key piece (bandos)" +
	"frozen key piece (zamorak)" +
	"frozen key piece (armadyl)"

type Drop struct {
	Name         string
	Rarity       string
	RawRarity    float64
	QuantityHigh int
	QuantityLow  int
	QuantityAvg  int
}

type DropTable struct {
	Drops []Drop
	Rolls int
}

type DataWrapper struct {
	Query struct {
		Results map[string]struct {
			Printouts struct {
				DropJSON []string `json:"Drop JSON"`
			} `json:"printouts"`
		} `json:"results"`
	} `json:"query"`
}

func main() {
	boss := "Kree'arra"
	u := GetQuery(boss)

	res, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(res.Body)

	dat := DataWrapper{}
	err2 := json.Unmarshal(b, &dat)
	if err2 != nil {
		log.Fatal(err2)
	}
	dt := dat.ParseDrops()
	x := dt.CheckTotalP()
	fmt.Printf("P = %.5f\n", x)
	r := dt.getRarest()
	fmt.Printf("Rarest = %d\n", r)

	dt.Sample(500)

}

func (d *DataWrapper) ParseDrops() DropTable {
	var output DropTable
	for _, v := range d.Query.Results {

		drop := Drop{}
		dropJSON := strings.Join(v.Printouts.DropJSON, "")
		dropJSON = strings.Trim(dropJSON, "{}")
		dropJSON = strings.ReplaceAll(dropJSON, ":", "")
		dropJSON = strings.ReplaceAll(dropJSON, ",", "")
		dropJSON = strings.ReplaceAll(dropJSON, "#", "")
		dropJSON = strings.ReplaceAll(dropJSON, "Alt Rarity", "")
		dropJSON = strings.ReplaceAll(dropJSON, "Name Notes", "")
		dropJSON = strings.ReplaceAll(dropJSON, "Rarity Notes", "")
		dropJSON = strings.ReplaceAll(dropJSON, " Dash", "")
		dropJSON = strings.ReplaceAll(dropJSON, " Dash", "")

		vals := strings.Split(dropJSON, "\"")
		var cleanedVals []string
		for i := 0; i < len(vals)-1; i++ {
			if len(vals[i]) > 0 {
				cleanedVals = append(cleanedVals, vals[i])
			}
		}
		for i := 0; i < len(cleanedVals)-1; i += 2 {
			switch cleanedVals[i] {
			case "Dropped item":
				drop.Name = cleanedVals[i+1]

			case "Rarity":
				drop.Rarity = cleanedVals[i+1]
				if cleanedVals[i+1] == "Always" {
					cleanedVals[i+1] = "1/1"
				}
				nd := strings.Split(cleanedVals[i+1], "/")

				num, err := strconv.ParseFloat(nd[0], 64)
				if err != nil {
					log.Fatal(err)
				}
				den, err := strconv.ParseFloat(nd[1], 64)
				if err != nil {
					log.Fatal(err)
				}
				drop.RawRarity = num / den

			case "Quantity High":
				drop.QuantityHigh, _ = strconv.Atoi(cleanedVals[i+1])
			case "Quantity Low":
				drop.QuantityLow, _ = strconv.Atoi(cleanedVals[i+1])
			}
		}
		drop.QuantityAvg = (drop.QuantityHigh + drop.QuantityLow) / 2

		if !strings.Contains(BannedItems, strings.ToLower(drop.Name)) {
			output.Drops = append(output.Drops, drop)
		}
		
	}
	return output
}

func GetQuery(boss string) url.URL {
	u := url.URL{
		Scheme:   "https",
		Host:     "oldschool.runescape.wiki",
		Path:     "api.php",
		RawQuery: "action=ask&format=json&query=",
	}
	q := fmt.Sprintf("[[-Has subobject::%s]]|[[Drop JSON::+]]|?Drop JSON", boss)
	u.RawQuery += url.QueryEscape(q)
	return u
}
