package dropsim

import (
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
	"net/url"
	"os"
	"strconv"
	"strings"
)

var SupportedBosses = []string{
	"General Graardor",
	"Commander Zilyana",
	"Kree'arra",
	"K'ril Tsutsaroth",
	"Dagannoth Rex",
	"Dagannoth Prime",
	"Dagannoth Supreme",
	"Giant Mole",
	"Sarachnis",
	"Corporeal Beast",
	"Kalphite Queen",
	"Barrows chest",
	"Zulrah",
}

// BannedItems Items that are too annoying to account for due to variable droprates
var BannedItems = "brimstone key" +
	"frozen key piece (saradomin)" +
	"frozen key piece (bandos)" +
	"frozen key piece (zamorak)" +
	"frozen key piece (armadyl)" +
	"Kq head (tattered)" +
	"Brimstone key"

// Tertiary TODO add separate tertiary rolls
var Tertiary = []string{
	"Clue scroll (hard)",
	"Clue scroll (elite)",
	"Long bone",
	"Curved bone",
	"Tattered moon page",
	"Tattered sun page",
	"Tattered temple page",
	"Grubby key",
	"Giant egg sac(full)",
}

// SpecialCases has items that drop together but the wiki api doesn't provide info on
var SpecialCases = map[string]map[string]string{
	"Kree'arra": {
		"Ranging potion (3)": "Super defence (3)",
	},
	"K'ril Tsutsaroth": {
		"Super attack (3)":  "Super strength (3)",
		"Super restore (3)": "Zamorak brew (3)",
	},
	"Commander Zilyana": {
		"Super defence (3)":  "Magic potion (3)",
		"Saradomin brew (3)": "Super restore (4)",
	},
}

type Drop struct {
	Name         string
	RawName      string `json:"Dropped item"`
	Rarity       string `json:"Rarity"`
	RawRarity    float64
	QuantityHigh int `json:"Quantity High"`
	QuantityLow  int `json:"Quantity Low"`
	QuantityAvg  int
	Rolls        int `json:"Rolls"`
	ImagePath    string
	ID           int
}

type DropList struct {
	Drops    []Drop
	Rolls    int
	BossName string
}

type DropWrapper struct {
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
	u := getQueryURL(boss)

	res, err := http.Get(u.String())
	if err != nil {
		panic(err)
	}
	b, err := io.ReadAll(res.Body)

	dat := DropWrapper{}
	err2 := json.Unmarshal(b, &dat)
	if err2 != nil {
		log.Fatal(err2)
	}
	dt := dat.ParseDrops(boss)
	x := dt.CheckTotalP()
	fmt.Printf("P = %.5f\n", x)
	r := dt.getRarest()
	fmt.Printf("Rarest = %d\n", r)

}

func (d *Drop) GetRawRarity() error {
	if d.Rarity == "Undefined" {
		d.RawRarity = 1
		return nil
	}
	ndc := strings.ReplaceAll(d.Rarity, ",", "")
	nd := strings.Split(ndc, "/")
	n, err := strconv.ParseFloat(nd[0], 64)
	if err != nil {
		log.Println(err)
		return err
	}
	de, err := strconv.ParseFloat(nd[1], 64)
	if err != nil {
		log.Fatal(err)
		return err
	}
	d.RawRarity = n / de
	return nil
}

func (d *DropWrapper) ParseDrops(boss string) DropList {
	var output DropList
	for _, v := range d.Query.Results {
		for _, d := range v.Printouts.DropJSON {
			drop := Drop{}
			if err := json.Unmarshal([]byte(d), &drop); err != nil {
				log.Fatalf("Error unmarshalling inner JSON: %v", err)
			}
			switch drop.Rarity {
			case "Always":
				drop.Rarity = "1/1"
			case "Common":
				drop.Rarity = "1/15"
			case "Uncommon":
				drop.Rarity = "1/40"
			case "Rare":
				drop.Rarity = "1/128"
			case "Once":
				drop.Rarity = "Undefined"
			}
			if err := drop.GetRawRarity(); err != nil {
				log.Fatalf("Error converting rarity: %v", err)
			}
			drop.Name = strings.ReplaceAll(drop.RawName, "#", " ")
			drop.QuantityAvg = (drop.QuantityHigh + drop.QuantityLow) / 2
			if !strings.Contains(BannedItems, strings.ToLower(drop.Name)) {
				output.Drops = append(output.Drops, drop)
			}
			err := drop.GetImage()
			if err != nil {
				log.Fatalf("Error getting image for %s: %v", drop.Name, err)
			}

		}

	}
	output.Rolls = output.Drops[0].Rolls
	output.BossName = boss
	return output
}

func GetAPIResponse(boss string) (DropList, error) {
	u := getQueryURL(boss)
	res, err := http.Get(u.String())
	if err != nil {
		return DropList{}, err
	}
	b, _ := io.ReadAll(res.Body)
	dat := DropWrapper{}
	err = json.Unmarshal(b, &dat)
	if err != nil {
		return DropList{}, err
	}
	dt := dat.ParseDrops(boss)
	return dt, nil
}

func getQueryURL(boss string) url.URL {
	u := url.URL{
		Scheme:   "https",
		Host:     "oldschool.runescape.wiki",
		Path:     "api.php",
		RawQuery: "action=ask&format=json&query=",
	}
	q := fmt.Sprintf("[[-Has subobject::%s]]|[[Drop JSON::+]]|?Drop JSON|limit=10000", boss)
	u.RawQuery += url.QueryEscape(q)
	return u
}

func (d *Drop) GetImage() error {
	path := fmt.Sprintf("images/%s.png", d.Name)
	if _, err := os.Stat(path); !os.IsNotExist(err) {
		return nil
	}
	u := url.URL{
		Scheme: "https",
		Host:   "oldschool.runescape.wiki",
		Path: "images/" +
			strings.ReplaceAll(strings.ReplaceAll(
				fmt.Sprintf("%s", d.RawName), " ", "_"), "#", ""),
	}
	l := strings.ToLower(d.Name)
	if strings.Contains(l, "arrow") ||
		strings.Contains(l, "bolt") ||
		strings.Contains(l, "seed") ||
		strings.Contains(l, "scales") ||
		strings.Contains(l, "brimstone key") {
		u.Path += url.PathEscape("_5")
	}
	if strings.Contains(l, "coins") {
		u.Path += url.PathEscape("_10000")
	}

	u.Path += url.PathEscape(".png")
	log.Println(u.String())
	r, err := http.Get(u.String())
	if err != nil {
		return err
	}
	defer r.Body.Close()
	data, err := io.ReadAll(r.Body)
	if err != nil {
		return err
	}
	os.WriteFile(path, data, 0666)
	d.ImagePath = fmt.Sprintf("images/%s.png", d.Name)
	return nil
}
