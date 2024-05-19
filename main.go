package main

import (
	"github.com/bwmarrin/discordgo"
	"log"
	"math/rand"
	"os"
	"os/signal"
	"strings"
	"syscall"
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

func main() {
	// get token
	Token, err := os.ReadFile(".env")
	if err != nil {
		log.Fatalf("Error reading token from .env file: %v", err)
	}
	if string(Token) == "" {
		log.Fatal("Token environment variable not set")
	}

	// create Discordgo session
	dg, err := discordgo.New("Bot " + string(Token))
	if err != nil {
		log.Fatalf("error creating Discord session: %v", err)
	}
	// add intents/permissions
	dg.Identify.Intents = discordgo.IntentsGuildMessages

	// register message responder function
	dg.AddHandler(messageCreate)

	err = dg.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}

	log.Println("Bot is now running.  Use ctrl-c to close.")
	sc := make(chan os.Signal, 1)
	signal.Notify(sc, syscall.SIGINT, syscall.SIGTERM, os.Interrupt)
	<-sc

	dg.Close()
}
