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

func FetchGuildEventCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	eventID := i.ApplicationCommandData().Options[0].StringValue()

	guildEvent, err := model.FetchGuildEvent(eventID)
	if err != nil {
		log.Printf("Error fetching event: %v", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching event data.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if guildEvent == nil || guildEvent.Id == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No event found for ID: %s", eventID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Build a large string for the event
	fullText := buildLargeEventString(guildEvent)

	// Chunk into ~1000-char pages
	pages := utils.ChunkString(fullText, 1000)
	if len(pages) == 0 {
		pages = []string{"No data available."}
	}

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	// Create and store the PaginationData
	utils.PaginationStore[paginationID] = &utils.PaginationData{
		Pages:     pages,
		PageIndex: 0,
		AuthorID:  i.Member.User.ID,
		Title:     "Fetched Guild Event",
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
	}
	utils.PaginationStore[paginationID] = data

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Failed to respond with guild event embed:", err)
	}
}
