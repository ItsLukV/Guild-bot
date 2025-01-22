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

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	fields := make([]utils.Section, 0)
	fields = append(fields, guildEvent)

	for _, data := range guildEvent.EventData {
		data.SetInLine(true)
		fields = append(fields, data)
	}

	// Create and store the PaginationData
	data := &utils.PaginationData{
		Fields:      fields,
		PageIndex:   0,
		Description: "",
		AuthorID:    i.Member.User.ID,
		Title:       "Fetched Guild Event",
		Footer:      fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		Color:       0x606060,
		CreatedAt:   time.Now(),
		PageSize:    5,
	}
	utils.PaginationStore[paginationID] = data

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Failed to respond with guild event embed:", err)
	}
}
