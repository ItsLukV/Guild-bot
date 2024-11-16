package commands

import (
	"fmt"
	"log"

	guildData "github.com/ItsLukV/Guild-bot/src/GuildData"
	"github.com/ItsLukV/Guild-bot/src/guildEvent"
	"github.com/bwmarrin/discordgo"
)

func createGuildEvent(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract options from the interaction
	options := i.ApplicationCommandData().Options

	var duration int
	var eventType guildEvent.EventType
	var description string

	// Map options by name for easy access
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Retrieve the duration (required, integer)
	if opt, ok := optionMap["duration"]; ok && opt.Type == discordgo.ApplicationCommandOptionInteger {
		duration = int(opt.IntValue())
	} else {
		respondWithError(s, i, "Invalid or missing 'duration' option.")
		return
	}

	// Retrieve the event type (required, string)
	if opt, ok := optionMap["event-type"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		eventType = guildEvent.EventType(opt.StringValue())
	} else {
		respondWithError(s, i, "Invalid or missing 'event-type' option.")
		return
	}

	// Retrieve the description (optional, string)
	if opt, ok := optionMap["description"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		description = opt.StringValue()
	} else {
		description = ""
	}

	// Create the event based on the event type
	event := guildEvent.NewGuildEvent(eventType, description, duration)

	// Add the event to the bot's events map
	g.Events[event.GetId()] = event

	// Respond to the user
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("Created a **%s** event with ID `%d`.\nDescription: %s", event.GetType(), event.GetId(), event.GetDescription()),
		},
	})
	log.Printf("Created a **%s** event with ID `%d`.\nDescription: %s", event.GetType(), event.GetId(), event.GetDescription())
	if err != nil {
		log.Printf("Failed to send interaction response: %v", err)
	}
}