package utils

import (
	"log"

	"github.com/bwmarrin/discordgo"
)

// Respond with a private message
func RespondWithErrorPrivate(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	embed := &discordgo.MessageEmbed{ // Use a pointer to match the expected type
		Title: message,
		Color: 0xFF0000,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed}, // Create a slice containing the embed
			Flags:  discordgo.MessageFlagsEphemeral,
		},
	})
	if err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}

// Respond with a public message
func RespondWithErrorPublic(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	embed := &discordgo.MessageEmbed{ // Use a pointer to match the expected type
		Title: message,
		Color: 0xFF0000,
	}

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{embed}, // Create a slice containing the embed
		},
	})
	if err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}
