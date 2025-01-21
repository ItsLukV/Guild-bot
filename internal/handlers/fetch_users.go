package handlers

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/model"
	"github.com/ItsLukV/Guild-bot/internal/restclient"
	"github.com/bwmarrin/discordgo"
)

func FetchUsersCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Build the target URL
	url := fmt.Sprintf("%s/api/users", config.GlobalConfig.ApiBaseURL)
	log.Println("Fetching users from:", url)
	// Fetch the user data
	data, err := restclient.FetchApi[model.UsersResponse](url)
	if err != nil {
		log.Println("Error fetching data:", err)
		// Respond with an ephemeral error message
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch data.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Handle the case of no users found
	if len(data.Users) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No users found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Create embed fields: one field per user
	fields := make([]*discordgo.MessageEmbedField, 0, len(data.Users))
	for _, user := range data.Users {
		username, err := model.FetchUsername(user.ActiveProfileUUID)
		if err != nil {
			log.Println("Error fetching username:", err)
			username = "Unknown"
		}
		value := fmt.Sprintf(
			"**User ID:** `%s`\n**Active Profile UUID:** `%s`\n**Fetch Data:** `%t`",
			user.ID,
			user.ActiveProfileUUID,
			user.FetchData,
		)

		fields = append(fields, &discordgo.MessageEmbedField{
			Name:  fmt.Sprintf("Username: %s", username),
			Value: value,
		})
	}

	// Build the embed
	embed := &discordgo.MessageEmbed{
		Title:       "Fetched Users",
		Description: fmt.Sprintf("We found **%d** user(s):", len(data.Users)),
		Color:       0x1F8B4C, // A greenish color; customize as desired
		Fields:      fields,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		},
		Timestamp: time.Now().UTC().Format(time.RFC3339), // Adds a timestamp
	}

	// Respond with the embed in the channel
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed},
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Println("Failed to respond with user embed:", err)
	}
}
