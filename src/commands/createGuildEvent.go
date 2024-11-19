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
		utils.RespondWithError(s, i, "Invalid or missing 'duration' option.")
		return
	}

	// Retrieve the event type (required, string)
	if opt, ok := optionMap["event-type"]; ok && opt.Type == discordgo.ApplicationCommandOptionInteger {
		eventType = guildData.EventType(opt.IntValue())
	} else {
		utils.RespondWithError(s, i, "Invalid or missing 'event-type' option.")
		return
	}

	// Retrieve the description (optional, string)
	if opt, ok := optionMap["description"]; ok && opt.Type == discordgo.ApplicationCommandOptionString {
		description = opt.StringValue()
	} else {
		description = ""
	}

	id, err := gonanoid.New()
	if err != nil {
		log.Printf("Failed to create a Id for guildEvent eventType: %v, description: %v , duration: %v", eventType, description, duration)
	}
	// Create the event based on the event type
	event := guildData.NewGuildEvent(id, eventType, description, time.Now(), duration)
	log.Print(len(g.Users))
	user := g.Users["365889587794673664"]
	event.AddUser(&user)
	// Add the event to the bot's events map
	g.Events[event.GetId()] = event

	db.GetInstance().SaveEvent(g.Events[event.GetId()])

	// Respond to the user
	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
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
