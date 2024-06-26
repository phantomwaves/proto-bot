package main

import (
	"context"
	"database/sql"
	"flag"
	"github.com/bwmarrin/discordgo"
	_ "github.com/mattn/go-sqlite3"
	"github.com/phantomwaves/proto/dropsim"
	"log"
	"os"
	"os/signal"
)

var (
	GuildID                = flag.String("guild", "524539398180700190", "guild ID - global mode active if empty")
	BotTokenPath           = flag.String("token", "", "bot token path")
	RemoveCommands         = flag.Bool("remove", false, "remove commands")
	DebugMode              = flag.Bool("debug", false, "debug mode - deleting database")
	minKills       float64 = 1
	//db             *sql.DB
)

var s *discordgo.Session

func getDropsimChoices(bosses []string) []*discordgo.ApplicationCommandOptionChoice {
	var output []*discordgo.ApplicationCommandOptionChoice
	for _, boss := range bosses {
		var choice = discordgo.ApplicationCommandOptionChoice{
			Name:  boss,
			Value: boss,
		}
		output = append(output, &choice)
	}
	return output
}

func init() {
	flag.Parse()
	var err error
	token := getToken(*BotTokenPath)
	s, err = discordgo.New("Bot " + token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	})
}

var commands = []*discordgo.ApplicationCommand{
	{
		Name:        "wiki",
		Description: "Wiki Search",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "search",
				Description: "Term to search the wiki for",
				Required:    true,
			},
		},
	},
	{
		Name:        "dropsim",
		Description: "Simulate the loot from n kills of a boss",
		Options: []*discordgo.ApplicationCommandOption{
			{
				Type:        discordgo.ApplicationCommandOptionString,
				Name:        "name",
				Description: "Full name of the boss to simulate",
				Required:    true,
				Choices:     getDropsimChoices(dropsim.SupportedBosses),
			},
			{
				Type:        discordgo.ApplicationCommandOptionInteger,
				Name:        "n",
				Description: "Number of kills",
				Required:    true,
				MaxValue:    10000,
				MinValue:    &minKills,
			},
		},
	},
}

func debugMode() {
	db, err := sql.Open("sqlite3", "./drops.db")
	if err != nil {
		log.Fatalf("Error opening db: %v", err)
	}
	defer db.Close()
	q1 := "PRAGMA writable_schema = 1;\n" +
		"delete from sqlite_master where type in ('table', 'index', 'trogger');\n" +
		"PRAGMA writable_schema = 0;\n" +
		"VACUUM;\n" +
		"PRAGMA INTEGRITY_CHECK;"
	db.ExecContext(context.Background(), q1)

}

func main() {
	if *DebugMode {
		debugMode()
	}
	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})
	err := s.Open()
	if err != nil {
		log.Fatalf("error opening connection: %v", err)
	}
	defer s.Close()

	//s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	//	if h, ok := CommandHandlers[i.ApplicationCommandData().Name]; ok {
	//		h(s, i)
	//	}
	//})
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands))

	for i, val := range commands {
		cmd, err2 := s.ApplicationCommandCreate(s.State.User.ID, "", val)
		if err2 != nil {
			log.Printf("cannot create command: %v\n", err)
		}
		registeredCommands[i] = cmd
	}

	s.Identify.Intents = discordgo.IntentsGuildMessages
	s.AddHandler(messageCreate)

	sc := make(chan os.Signal, 1)
	signal.Notify(sc, os.Interrupt)
	log.Println("Press Ctrl-c to exit.")
	<-sc

	registeredCommands, err = s.ApplicationCommands(s.State.User.ID, *GuildID)
	if err != nil {
		log.Fatalf("Could not fetch registered commands: %v", err)
	}

	if *RemoveCommands {
		log.Println("Removing commands...")
		for _, cmd := range registeredCommands {
			err2 := s.ApplicationCommandDelete(s.State.User.ID, *GuildID, cmd.ID)
			if err2 != nil {
				log.Printf("cannot delete command: %v\n", err2)
			}
		}
	}
	log.Println("shutting down.")

}

func getToken(path string) string {
	if path == "" {
		path = ".env"
	}
	Token, err := os.ReadFile(path)
	if err != nil {
		log.Fatalf("error reading token from .env file: %v", err)
	}
	if string(Token) == "" {
		log.Fatal("token environment variable not set")
	}
	return string(Token)
}
