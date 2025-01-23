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

func FetchUsersCommand(s *discordgo.Session, i *discordgo.InteractionCreate, pm *utils.PaginatedSessions) {
	// Fetch the user data
	users, err := model.FetchUsers()
	if err != nil {
		log.Println("Error fetching users:", err)
		// Respond with an ephemeral error message
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch users.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Handle the case of no users found
	if len(users) == 0 {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No users found.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var fields []utils.Section
	for i := range users {
		fields = append(fields, &users[i])
	}

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	// Create and store the PaginationData
	data := &utils.PaginationData{
		Fields:    fields,
		PageIndex: 0,
		AuthorID:  i.Member.User.ID,
		Title:     "Fetched Users",
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
		PageSize:  5,
	}
	pm.Put(paginationID, data)

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Failed to respond with paginated users:", err)
	}
}
