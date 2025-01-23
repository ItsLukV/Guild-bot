package autocompletions

import (
	"github.com/ItsLukV/Guild-bot/internal/model"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func GuildEventAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Sanity check: ensure we have the right command & option
	data := i.ApplicationCommandData()
	if data.Name != "fetch_guild_event" {
		return
	}
	if len(data.Options) == 0 {
		return
	}

	// The user’s typed input is in i.Data.Options[0].Value (or i.ApplicationCommandData().Options[0].StringValue())
	userTyped := data.Options[0].StringValue()

	// Fetch all events
	events, err := model.FetchGuildEvents()
	if err != nil {
		log.Println("Error fetching guild events:", err)
		// If something fails, we could just respond with an empty list or a generic error.
		respondAutocomplete(s, i, nil)
		return
	}

	// Filter or rank events based on `userTyped`
	// For example: show only events whose ID *contains* the typed text
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, event := range events {
		if strings.Contains(strings.ToLower(event.Id), strings.ToLower(userTyped)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  event.Id, // what the user sees in the dropdown
				Value: event.Id, // the actual value sent if selected
			})
		}
	}

	// If user typed nothing (or if you want to show all), you can show everything
	if userTyped == "" && len(choices) == 0 {
		for _, event := range events {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  event.Id,
				Value: event.Id,
			})
		}
	}

	respondAutocomplete(s, i, choices)
}

func respondAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate, choices []*discordgo.ApplicationCommandOptionChoice) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionApplicationCommandAutocompleteResult,
		Data: &discordgo.InteractionResponseData{
			Choices: choices,
		},
	})
	if err != nil {
		log.Println("Failed to send autocomplete response:", err)
	}
}
