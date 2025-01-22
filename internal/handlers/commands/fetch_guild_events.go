package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/model"
	"github.com/ItsLukV/Guild-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func FetchGuildEventsCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Fetch the guild events
	events, err := model.FetchGuildEvents()
	if err != nil {
		log.Println("Error fetching guild events:", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch guild events.",
				Flags:   discordgo.MessageFlagsEphemeral, // Only the user sees this
			},
		})
		return
	}

	// Handle the case of no events found
	if len(events) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No guild events found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var fields []utils.Section
	for i := range events {
		fields = append(fields, &events[i])
	}

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	// Create and store the PaginationData
	data := &utils.PaginationData{
		Fields:    fields,
		PageIndex: 0,
		AuthorID:  i.Member.User.ID,
		Title:     "Fetched Guild Events",
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
		PageSize:  5,
	}
	utils.PaginationStore[paginationID] = data

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Failed to respond with guild events embed:", err)
	}
}

// buildAllEventsString returns a large string describing each event
func buildAllEventsString(events []model.GuildEvent) string {
	var output string

	output += fmt.Sprintf("We found **%d** event(s):\n\n", len(events))

	for idx, event := range events {
		// Convert user IDs to names
		userNames, err := event.GetUserNames()
		if err != nil {
			log.Printf("Error getting user names: %v\n", err)
			userNames = event.Users
		}

		eventText := fmt.Sprintf(
			"[%d] Event ID: %s\n"+
				"Users: %v\n"+
				"Start Time: %s\n"+
				"Duration: %dh\n"+
				"Type: %s\n"+
				"Hidden: %t\n\n",
			idx+1,
			event.Id,
			userNames,
			event.StartTime.Format(time.RFC1123),
			event.Duration,
			event.Type,
			event.IsHidden,
		)

		output += eventText
	}

	return output
}
