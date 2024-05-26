package main

import (
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

		u := dropsim.GetQuery(options[0].StringValue())
		res, err := http.Get(u.String())
		if err != nil {
			log.Printf("wiki request failed. %v\n", err)
		}
		b, _ := io.ReadAll(res.Body)
		dat := dropsim.DataWrapper{}
		err = json.Unmarshal(b, &dat)
		if err != nil {
			log.Printf("error unmarshaling json: %v\n", err)
		}
		dt := dat.ParseDrops()
		//log.Printf("Simulating loot from %v %v kills...\n",
		//	options[1].IntValue(), options[0].StringValue())
		response := fmt.Sprintf("Simulating loot from %v %v kills...\n",
			options[1].IntValue(), options[0].StringValue())
		response += dt.Sample(int(options[1].IntValue()))

		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: response,
			},
		})

	},
}