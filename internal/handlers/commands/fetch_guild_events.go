package commands

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

func FetchGuildEventsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	log.Println("Fetching guild events")
}
