package commands

import (
	"fmt"
	"github.com/ItsLukV/Guild-bot/internal/model"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func FetchGuildEventsCommand(s *discordgo.Session, i *discordgo.InteractionCreate, pm *utils.PaginatedSessions) {
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
	pm.Put(paginationID, data)

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Failed to respond with guild events embed:", err)
	}
}
