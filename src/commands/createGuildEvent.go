package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/src/db"
	guildData "github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
	gonanoid "github.com/matoous/go-nanoid/v2"
)

func createGuildEvent(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Extract options from the interaction
	options := i.ApplicationCommandData().Options

	var duration int
	var eventType guildData.EventType
	var eventName string
	var description string
	var startDate time.Time
	var hidden bool

	// Map options by name for easy access
	optionMap := make(map[string]*discordgo.ApplicationCommandInteractionDataOption)
	for _, opt := range options {
		optionMap[opt.Name] = opt
	}

	// Retrieve the duration (required, integer)
	if opt, ok := optionMap["duration"]; ok && opt.Type == discordgo.ApplicationCommandOptionInteger {
		duration = int(opt.IntValue())
	} else {
		utils.RespondWithErrorPrivate(s, i, "Invalid or missing 'duration' option.")
		return
	}

	// Retrieve the event type (required, string)
	if opt, ok := optionMap["event-type"]; ok && opt.Type == discordgo.ApplicationCommandOptionInteger {
		eventType = guildData.EventType(opt.IntValue())
	} else {
		utils.RespondWithErrorPrivate(s, i, "Invalid or missing 'event-type' option.")
		return
	}

	// Retrieve the description (optional, string)
	if opt, ok := optionMap["description"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		description = opt.StringValue()
	} else {
		description = ""
	}

	// Retrieve start date (optional, string)
	if opt, ok := optionMap["start-date"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		const layout = "15:04 02-01" // Define the custom timestamp layout

		// Parse the start date
		startDateStr := opt.StringValue()
		log.Println(startDateStr)
		var err error
		startDate, err = time.Parse(layout, startDateStr)
		if err != nil {
			utils.RespondWithErrorPrivate(s, i, "Invalid timestamp format. Please use 'HH:mm DD-MM'.")
			return
		}

		// If parsing succeeds, proceed with startDate
	} else {
		// If the option is not provided, set startDate to zero value
		startDate = time.Time{}
	}

	// Retrieve the eventName (optional, string)
	if opt, ok := optionMap["event-name"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		eventName = opt.StringValue()
	} else {
		eventName = ""
	}

	// Retrieve the hidden tag (optional, string)
	if opt, ok := optionMap["hidden"]; ok && opt.Type == discordgo.ApplicationCommandOptionBoolean {
		hidden = opt.BoolValue()
	} else {
		hidden = false
	}
	alphabet := "abcdefghijklmnopqrstuvwxyz1234567890ABCDEFGHIJKLMNOPQRSTUVWXYZ"
	id, err := gonanoid.Generate(alphabet, 21)
	if err != nil {
		log.Printf("Failed to create a Id for guildEvent eventType: %v, description: %v , duration: %v", eventType, description, duration)
	}

	// Create the event based on the event type
	event := guildData.NewGuildEvent(id, eventName, eventType, description, startDate, time.Now().UTC(), duration, false, hidden)
	// Add the event to the bot's events map
	g.Events[event.GetId()] = event

	db.GetInstance().SaveEvent(g.Events[event.GetId()])

	// Respond to the user

	// embed := discordgo.MessageEmbed{
	// 	Title: message,
	// 	Color: 0xFF0000,
	// }

	field := []*discordgo.MessageEmbedField{
		{
			Name:   "Duration",
			Value:  fmt.Sprintf("%d hours", event.GetDuration()),
			Inline: true,
		},
	}

	if !event.GetStartTime().IsZero() {
		outputStartDate := fmt.Sprintf("<t:%v:R>", event.GetStartTime().Unix())
		field = append(field, &discordgo.MessageEmbedField{
			Name:   "Start Date",
			Value:  outputStartDate, // Adjust format as needed
			Inline: true,
		})
	}
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds: []*discordgo.MessageEmbed{
				{
					Title:       fmt.Sprintf("New Event: %s", event.GetType()),
					Description: fmt.Sprintf("**Event ID:** `%s`\n**Event Name:**`%s`\n**Description:** %s", event.GetId(), event.GetEventName(), event.GetDescription()),
					Color:       0x1F8B4C, // A greenish color
					Fields:      field,
					Footer: &discordgo.MessageEmbedFooter{
						Text: "Click the button below to join the event.",
					},
					Timestamp: time.Now().UTC().Format(time.RFC3339), // Embed timestamp
				},
			},
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Emoji: &discordgo.ComponentEmoji{
								Name: "📜",
							},
							Label:    "Join Event",
							Style:    discordgo.PrimaryButton,
							CustomID: fmt.Sprintf("joinEvent_%s", event.GetId()),
						},
					},
				},
			},
		},
	})
	log.Printf("Created a **%s** event with ID `%s`.\nDescription: %s", event.GetType(), event.GetId(), event.GetDescription())
	if err != nil {
		log.Printf("Failed to send interaction response: %v", err)
	}
}
