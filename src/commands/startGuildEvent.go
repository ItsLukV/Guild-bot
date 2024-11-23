package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/src/db"
	"github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/utils"

	"github.com/bwmarrin/discordgo"
)

func startguildEvent(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	switch i.Type {
	// Command submition
	case discordgo.InteractionApplicationCommand:
		// Handle the command execution when the user submits the command
		data := i.ApplicationCommandData()

		// Retrieve the event ID from the options
		var eventID string
		if len(data.Options) > 0 && data.Options[0].Type == discordgo.ApplicationCommandOptionString {
			eventID = string(data.Options[0].StringValue())
		} else {
			utils.RespondWithErrorPrivate(s, i, "Invalid or missing 'event-id' option.")
			return
		}

		// Find the event in your bot's event map
		event, exists := g.Events[eventID]
		if !exists {
			utils.RespondWithErrorPrivate(s, i, fmt.Sprintf("Event with ID %v not found.", eventID))
			return
		}

		// Check if there are any users
		if len(event.GetUsers()) == 0 {
			utils.RespondWithErrorPrivate(s, i, "Event has no users.")
			return
		}

		// Start the event
		err := event.Start()
		if err != nil {
			utils.RespondWithErrorPrivate(s, i, fmt.Sprintf("Failed to start event: %v", err))
			return
		}
		db.GetInstance().SaveEventData(event)
		db.GetInstance().UpdateEvent(event)

		// Create an embed to respond to the user
		embed := &discordgo.MessageEmbed{
			Title:       fmt.Sprintf("Event '%s' Started", event.GetDescription()),
			Description: fmt.Sprintf("Event with ID: %s has successfully started.", event.GetId()),
			Color:       0x00FF00, // Green color to indicate success
			Fields: []*discordgo.MessageEmbedField{
				{
					Name:   "Event Type",
					Value:  event.GetType().String(),
					Inline: true,
				},
				{
					Name:   "Users Joined",
					Value:  fmt.Sprintf("%d", len(event.GetUsers())),
					Inline: true,
				},
				{
					Name:   "Users Joined",
					Value:  fmt.Sprintf("Start Date: <t:%v:R>, End Date: <t:%v:R>", event.GetStartTime().Unix(), time.Now().UTC().AddDate(0, 0, event.GetDuration()/24)),
					Inline: true,
				},
			},
			Footer: &discordgo.MessageEmbedFooter{
				Text: "Event started successfully!",
			},
		}

		// Respond with the embed
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Embeds: []*discordgo.MessageEmbed{embed},
			},
		})
		if err != nil {
			log.Printf("Failed to send interaction response: %v", err)
		}
	// Autocomplete options introduce a new interaction type (8) for returning custom autocomplete results.
	case discordgo.InteractionApplicationCommandAutocomplete:
		// Collect event IDs of events that have not been started
		var choices []*discordgo.ApplicationCommandOptionChoice
		for _, event := range g.Events {
			if !event.GetIsActive() || !event.HasEnded() {
				choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
					Name:  fmt.Sprintf("Event %v: %s", event.GetId(), event.GetType()),
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
