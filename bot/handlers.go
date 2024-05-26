package main

import (
	"github.com/bwmarrin/discordgo"
	"math/rand"
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
}
