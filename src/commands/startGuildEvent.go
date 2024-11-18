package commands

import (
	"fmt"
	"log"

	guildData "github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/bwmarrin/discordgo"
)

func startguildEvent(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	// Command submition
	case discordgo.InteractionApplicationCommand:
		// Handle the command execution when the user submits the command
		data := i.ApplicationCommandData()

		// Retrieve the event ID from the options
		var eventID int
		if len(data.Options) > 0 && data.Options[0].Type == discordgo.ApplicationCommandOptionInteger {
			eventID = int(data.Options[0].IntValue())
		} else {
			respondWithError(s, i, "Invalid or missing 'event-id' option.")
			return
		}

		// Find the event in your bot's event map
		event, exists := g.Events[eventID]
		if !exists {
			respondWithError(s, i, fmt.Sprintf("Event with ID %d not found.", eventID))
			return
		}

		// Start the event
		err := event.Start()
		if err != nil {
			respondWithError(s, i, fmt.Sprintf("Failed to start event: %v", err))
			return
		}

		// Respond to the user confirming the event has started
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("Event '%s' (ID: %d) has been started.", event.GetDescription(), event.GetId()),
			},
		})
		if err != nil {
			log.Printf("Failed to send interaction response: %v", err)
		} else {
			event.Start()
		}
	// Autocomplete options introduce a new interaction type (8) for returning custom autocomplete results.
	case discordgo.InteractionApplicationCommandAutocomplete:
		// Collect event IDs of events that have not been started
		var choices []*discordgo.ApplicationCommandOptionChoice
		for _, event := range g.Events {
			if !event.GetIsActive() {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  fmt.Sprintf("Event %d: %s", event.GetId(), event.GetType()),
					Value: event.GetId(),
				})
			}
		}

		// Limit to 25 choices as per Discord's limit
		if len(choices) > 25 {
			choices = choices[:25]
		}

		// Respond with the choices
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionApplicationCommandAutocompleteResult,
			Data: &discordgo.InteractionResponseData{
				Choices: choices,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to Autocomplete interaction: %v", err)
		}
	}
}