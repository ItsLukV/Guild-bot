package commands

import (
	"fmt"
	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/model"
	"github.com/ItsLukV/Guild-bot/internal/utils"
	"log"
	"time"

	"github.com/bwmarrin/discordgo"
)

func FetchUserCommand(s *discordgo.Session, i *discordgo.InteractionCreate, pm *utils.PaginatedSessions) {
	// Parse the command options
	userID := i.ApplicationCommandData().Options[0].StringValue()

	// Fetch the user data
	fullUserData, err := model.FetchUser(userID)
	if err != nil {
		log.Println("Error fetching specific user:", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch specific user.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if fullUserData == nil || fullUserData.User.Id == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "No user found for that ID.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	var fields []utils.Section

	fields = append(fields, &fullUserData.User)
	fields = append(fields, &fullUserData.DianaData)
	fields = append(fields, &fullUserData.DungeonsData)

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	// Create and store the PaginationData
	minecraftName, err := model.FetchUsername(userID)
	if err != nil {
		log.Println("Error fetching username:", err)
		minecraftName = userID
	}

	data := &utils.PaginationData{
		Fields:    fields,
		PageIndex: 0,
		AuthorID:  i.Member.User.ID,
		Title:     fmt.Sprintf("Fetched User: %s", minecraftName),
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
		PageSize:  5,
	}
	pm.Put(paginationID, data)

	if err := utils.SendInitialPaginationResponse(s, i, paginationID, data); err != nil {
		log.Println("Error sending initial pagination response:", err)
	}
}
