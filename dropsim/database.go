package dropsim

import (
	"context"
	"database/sql"
	"errors"
	"fmt"
	"log"
	"strings"
)

// GetDropsTable retrieves the table for a boss from db if it already exists
func GetDropsTable(db *sql.DB, boss string) ([]Drop, error) {
	var drops []Drop
	boss = strings.ReplaceAll(boss, " ", "_")
	boss = strings.ReplaceAll(boss, "'", "_")
	q1 := fmt.Sprintf("SELECT * FROM %s_tbl", boss)
	rows, err := db.QueryContext(context.Background(), q1)
	if err != nil {
		return nil, err
	}
	defer rows.Close()
	for rows.Next() {
		var drop Drop
		if err = rows.Scan(&drop.ID, &drop.Name, &drop.QuantityAvg); err != nil {
			return nil, err
		}
		drops = append(drops, drop)
	}
	return drops, nil
}

// AddDropsTable adds the table for a boss to db
func AddDropsTable(db *sql.DB, drops []Drop, boss string) error {
	tableName := strings.ReplaceAll(boss, " ", "_")
	tableName = strings.ReplaceAll(tableName, "'", "_")
	tableName += "_tbl"
	q1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, item_name TEXT NOT NULL, quantity INTEGER NOT NULL)", tableName)
	_, err := db.ExecContext(context.Background(), q1)
	if err != nil {
		return err
	}
	q2 := fmt.Sprintf("INSERT INTO %s (item_name, quantity)\nVALUES\n", tableName)
	for _, drop := range drops {
		q2 += fmt.Sprintf("(\"%s\", %d),\n", drop.Name, drop.QuantityAvg)
	}
	q2 = strings.TrimRight(q2, ",\n")
	_, err = db.ExecContext(context.Background(), q2)
	if err != nil {
		return err
	}
	return nil
}

// AddBoss adds boss drops info to db
func AddBoss(db *sql.DB, drops DropList, boss string) error {
	tableName := strings.ReplaceAll(boss, " ", "_")
	tableName = strings.ReplaceAll(tableName, "'", "_")

	q1 := fmt.Sprintf("CREATE TABLE IF NOT EXISTS %s (id INTEGER PRIMARY KEY AUTOINCREMENT, item_name TEXT NOT NULL, quantity INTEGER NOT NULL, rarity TEXT NOT NULL, rolls INT NOT NULL)", tableName)
	_, err := db.ExecContext(context.Background(), q1)
	if err != nil {
		return err
	}
	q2 := fmt.Sprintf("INSERT INTO %s (item_name, quantity, rarity, rolls)\nVALUES\n", tableName)
	for _, drop := range drops.Drops {
		q2 += fmt.Sprintf("(\"%s\", %d, \"%s\", %d),\n", drop.Name, drop.QuantityAvg, drop.Rarity, drop.Rolls)
	}
	q2 = strings.TrimRight(q2, ",\n")
	log.Println(q2)
	_, err = db.ExecContext(context.Background(), q2)
	if err != nil {
		return err
	}
	return nil
}

// GetBoss retrieves boss drops info from db
func GetBoss(db *sql.DB, boss string) (DropList, error) {
	tableName := strings.ReplaceAll(boss, " ", "_")
	tableName = strings.ReplaceAll(tableName, "'", "_")
	output := DropList{}
	q1 := fmt.Sprintf("SELECT * FROM %s", tableName)
	rows, err := db.QueryContext(context.Background(), q1)
	if err != nil {
		return DropList{}, err
	}
	defer rows.Close()
	for rows.Next() {
		var drop Drop
		if err = rows.Scan(&drop.ID, &drop.Name, &drop.QuantityAvg, &drop.Rarity, &drop.Rolls); err != nil {
			return DropList{}, err
		}
		drop.GetRawRarity()
		output.Drops = append(output.Drops, drop)
	}
	if len(output.Drops) == 0 {
		return DropList{}, errors.New("no drops found")
	}
	output.Rolls = output.Drops[0].Rolls
	output.BossName = boss
	return output, nil
}
