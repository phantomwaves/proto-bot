package main

import (
	"bytes"
	"database/sql"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/phantomwaves/proto/dropsim"
	"log"
	"math/rand"
	"strings"
)

func validBoss(boss string) bool {
	for _, b := range dropsim.SupportedBosses {
		if b == boss {
			return true
		}
	}
	return false
}

func messageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {

	// ignore messages from self
	if m.Author.ID == s.State.User.ID {
		return
	}

	if strings.Contains(m.Content, "73") {
		s.ChannelMessageSend(m.ChannelID, "73 loooooool :joy::joy::joy:")
	}
	if strings.Contains(strings.ToLower(m.Content), "ooo") {
		x, y := 2, 20
		n := rand.Intn(y) + x

		msg := "O" + strings.Repeat("o", n)

		s.ChannelMessageSend(m.ChannelID, msg)
	}
}

var CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){

	"wiki": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		options := i.ApplicationCommandData().Options
		url := NewURL(options[0].StringValue())

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: url.String(),
			},
		})
	},

	"dropsim": func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		log.Println("received request for dropsim")
		options := i.ApplicationCommandData().Options
		boss := options[0].StringValue()
		n := options[1].IntValue()
		if !validBoss(boss) {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Sorry, %s is not supported yet.", i.ApplicationCommandData().Options[0].StringValue()),
				},
			})
		} else {
			db, err := sql.Open("sqlite3", "./drops.db")
			if err != nil {
				log.Fatalf("Error opening db: %v", err)
			}
			defer db.Close()

			dt, err := dropsim.GetBoss(db, boss)
			if err != nil {
				log.Printf("Error reading boss: %v", err)
			}
			if len(dt.Drops) == 0 {
				log.Printf("Database entry doesn't exist. Sending API query for %v %v...", n, boss)
				dt, err = dropsim.GetAPIResponse(boss)
				if err != nil {
					log.Fatalf("Error making api query: %v", err)
				}
				err := dropsim.AddBoss(db, dt, boss)
				if err != nil {
					log.Fatalf("Error adding boss to DB: %v", err)
				}
			}

			tbl, err := dropsim.GetDropsTable(db, boss)
			if err != nil {
				log.Printf("Error reading DB: %v", err)
			}
			if len(tbl) == 0 {
				tbl = dt.MakeDropTable()
				err = dropsim.AddDropsTable(db, tbl, boss)
				if err != nil {
					log.Fatalf("Error adding drops table to DB: %v", err)
				}
			}
			log.Printf("Sampling %v %v...", n, boss)
			itemCounts := dt.Sample(int(n), tbl, boss)

			r := dropsim.ResponseImage{
				Title: fmt.Sprintf("Loot from %v %v", n, boss),
			}
			log.Printf("Making image %v %v...", n, boss)
			r.MakeResponse(itemCounts)
			img, err := r.GetScreenshot(r.Filepath)
			if err != nil {
				log.Printf("error getting image: %v\n", err)
			}
			rd := bytes.NewReader(img)

			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Files: []*discordgo.File{
						{
							ContentType: "image",
							Name:        "Drops.png",
							Reader:      rd,
						},
					},
				},
			})
		}

	},
}
