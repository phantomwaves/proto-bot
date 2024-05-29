package main

import (
	"bytes"
	"encoding/json"
	"fmt"
	"github.com/bwmarrin/discordgo"
	"github.com/phantomwaves/proto/dropsim"
	"io"
	"log"
	"math/rand"
	"net/http"
	"strings"
)

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
		options := i.ApplicationCommandData().Options
		if !func() bool {
			for _, b := range dropsim.SupportedBosses {
				if b == i.ApplicationCommandData().Options[0].StringValue() {
					return true
				}
			}
			return false
		}() {
			s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: fmt.Sprintf("Sorry, %s is not supported yet.", i.ApplicationCommandData().Options[0].StringValue()),
				},
			})
		} else {
			log.Printf("Making query for %v %v...", options[1].IntValue(), options[0].StringValue())
			u := dropsim.GetDropsData(options[0].StringValue())
			log.Printf("Sending query for %v %v...", options[1].IntValue(), options[0].StringValue())
			res, err := http.Get(u.String())
			if err != nil {
				log.Printf("request failed. %v\n", err)
			}
			b, _ := io.ReadAll(res.Body)
			dat := dropsim.DropWrapper{}
			err = json.Unmarshal(b, &dat)
			if err != nil {
				log.Printf("error unmarshalling json: %v\n", err)
			}
			dt := dat.ParseDrops()
			log.Printf("Sampling %v %v...", options[1].IntValue(), options[0].StringValue())
			itemCounts := dt.Sample(int(options[1].IntValue()))

			r := dropsim.ResponseImage{
				Title: fmt.Sprintf("Loot from %v %v", options[1].IntValue(), options[0].StringValue()),
			}
			log.Printf("Making image %v %v...", options[1].IntValue(), options[0].StringValue())
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
