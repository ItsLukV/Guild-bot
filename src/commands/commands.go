package commands

import (
	"log"
	"strings"

	guildData "github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/guildEvent"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "register",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Link your minecraft account with your discord account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "minecraft-username",
					Description: "Minecraft username to be added",
					Required:    true,
				},
			},
		},
		{
			Name: "create-event",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Create a guild event.",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "duration",
					Description: "The duration of the guild event in hours",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionInteger,
					Name:        "event-type",
					Description: "The type of the guild event",
					Required:    true,
					Choices: generateChoices([]guildEvent.EventType{
						guildEvent.Slayer,
						guildEvent.Diana,
						guildEvent.Dungeons,
					}),
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "description",
					Description: "Add a description for the guild event",
					Required:    false,
				},
			},
		},
		{
			Name: "start-event",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Create",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionInteger,
					Name:         "event-id",
					Description:  "The type of the guild event",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		"register":     registerAccount,
		"create-event": createGuildEvent,
		"start-event":  startguildEvent,
	}
)

func generateChoices(events []guildEvent.EventType) []*discordgo.ApplicationCommandOptionChoice {
	choices := make([]*discordgo.ApplicationCommandOptionChoice, len(events))
	for i, event := range events {
		choices[i] = &discordgo.ApplicationCommandOptionChoice{
			Name:  strings.ToLower(event.String()),
			Value: event,
		}
	}
	return choices
}

func respondWithError(s *discordgo.Session, i *discordgo.InteractionCreate, message string) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: message,
			Flags:   discordgo.MessageFlagsEphemeral, // Only visible to the user
		},
	})
	if err != nil {
		log.Printf("Failed to send error response: %v", err)
	}
}