package autocompletions

import (
	"github.com/ItsLukV/Guild-bot/internal/model"
	"log"
	"strings"

	"github.com/bwmarrin/discordgo"
)

func UserAutocomplete(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Sanity check: ensure we have the right command & option
	data := i.ApplicationCommandData()
	if data.Name != "fetch_user" {
		return
	}
	if len(data.Options) == 0 {
		return
	}

	// The user’s typed input is in i.Data.Options[0].Value (or i.ApplicationCommandData().Options[0].StringValue())
	userTyped := data.Options[0].StringValue()

	// Fetch all users
	users, err := model.FetchUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		respondAutocomplete(s, i, nil)
		return
	}

	// Filter or rank events based on `user name`
	var choices []*discordgo.ApplicationCommandOptionChoice
	for _, user := range users {
		name, err := model.FetchUsername(user.Id)
		if err != nil {
			log.Println("Error fetching username:", err)
			name = user.Id
		}

		if strings.Contains(strings.ToLower(name), strings.ToLower(userTyped)) {
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: user.Id,
			})
		}
	}

	// If user typed nothing (or if you want to show all), you can show everything
	if userTyped == "" && len(choices) == 0 {
		for _, user := range users {
			name, err := model.FetchUsername(user.Id)
			if err != nil {
				log.Println("Error fetching username:", err)
				name = user.Id
			}
			choices = append(choices, &discordgo.ApplicationCommandOptionChoice{
				Name:  name,
				Value: user.Id,
			})
		}
	}

	respondAutocomplete(s, i, choices)
}
