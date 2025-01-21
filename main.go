package main

import (
	"flag"
	"log"
	"os"
	"os/signal"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/discord"
	"github.com/joho/godotenv"
)

var RemoveCommands = flag.Bool("rmcmd", true, "Remove all commands after shutdown or not")

func init() {
	flag.Parse()
}

func main() {
	log.SetFlags(log.LstdFlags | log.Lshortfile)
	// Load environment variables
	if err := godotenv.Load(".env"); err != nil {
		log.Println("Error loading .env file")
	} else {
		log.Println(".env file loaded successfully")
	}
	// Load the configuration
	config.LoadConfig()

	// Initialize Discord service
	discordSvc := discord.New(config.GlobalConfig.BotToken)

	// Add command handlers
	discordSvc.AddCommandHandlers()

	// Start the bot
	if err := discordSvc.Start(); err != nil {
		log.Fatalf("Cannot open the session: %v", err)
	}

	// Register commands
	registeredCommands := discordSvc.RegisterCommands()

	// Wait for a signal to gracefully shutdown the bot
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, os.Interrupt)
	log.Println("Press Ctrl+C to exit")
	<-stop

	if *RemoveCommands {
		log.Println("Removing commands...")
		discordSvc.RemoveCommands(registeredCommands)
	}

	log.Println("Gracefully shutting down.")
}
