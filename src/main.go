package main

import (
	"flag"
	"log"
	"os"
	"os/signal"
	"time"

	"github.com/ItsLukV/Guild-bot/src/commands"
	"github.com/ItsLukV/Guild-bot/src/db"
	guildData "github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/bwmarrin/discordgo"
)

var (
	discord_token  = os.Getenv("DISCORD_TOKEN")
	GuildID        = os.Getenv("GUILD_ID")
	RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdowning or not")
)

var s *discordgo.Session

var data guildData.GuildBot

func init() {
	flag.Parse()
	var err error
	s, err = discordgo.New("Bot " + discord_token)
	if err != nil {
		log.Fatalf("Invalid bot parameters: %v", err)
	}
	// Loading data from db
	data = guildData.GuildBot{
		DiscordSession: s,
	}
	err = db.LoadGuildBot(&data)
	if err != nil {
		log.Println("Failed to load guildBot from the Database")
		log.Panicf("Database connection error: %v", err)
	}

	// Handlers
	s.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {

		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := commands.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(&data, s, i)
			}
		case discordgo.InteractionMessageComponent:
			commands.HandleButtonInteraction(&data, s, i)
		// case discordgo.InteractionModalSubmit:
		// case discordgo.InteractionPing:
		default:
			log.Panicf("unexpected discordgo.InteractionType: %#v", i.Interaction.Type)
		}
	})
}

func main() {
	// Create a ticker that ticks at the specified interval
	ticker := time.NewTicker(3 * time.Hour)
	data.StartEventUpdater(*ticker)
	defer ticker.Stop()

	s.AddHandler(func(s *discordgo.Session, r *discordgo.Ready) {
		log.Printf("Logged in as: %v#%v", s.State.User.Username, s.State.User.Discriminator)
	})

	err := s.Open()
	if err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	log.Println("Adding commands...")
	registeredCommands := make([]*discordgo.ApplicationCommand, len(commands.Commands))
	for i, v := range commands.Commands {
		cmd, err := s.ApplicationCommandCreate(s.State.User.ID, GuildID, v)
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		registeredCommands[i] = cmd
	}

	defer s.Close()

	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")

		for _, v := range registeredCommands {
			err := s.ApplicationCommandDelete(s.State.User.ID, GuildID, v.ID)
			if err != nil {
				log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
			}
		}
	}

	db.GetInstance().Close()
	log.Println("Gracefully shutting down.")
}
