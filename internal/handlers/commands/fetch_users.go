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

func FetchUsersCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Fetch the user data
	users, err := model.FetchUsers(config.GlobalConfig.ApiBaseURL)
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

	// Build a large string describing all users
	fullText := buildLargeUsersString(users)

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
		Title:     "Fetched Users",
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
	}

	// Create an embed for page 0
	embed := utils.MakePaginationEmbed(utils.PaginationStore[paginationID])

	// Build pagination buttons
	components := utils.MakePaginationComponents(paginationID, 0, len(pages))

	// Send the initial response
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	}); err != nil {
		log.Println("Failed to respond with paginated users:", err)
	}
}

func buildLargeUsersString(users []model.User) string {
	var output string

	output += fmt.Sprintf("We found **%d** user(s):\n\n", len(users))

	for _, user := range users {
		// Get the username for this user
		username, err := model.FetchUsername(user.ID)
		if err != nil {
			log.Println("Error fetching username:", err)
			username = "Unknown"
		}

		output += fmt.Sprintf(
			"**<@%s>**\n"+
				"**Minecraft username:** `%s`\n"+
				"**Fetching Data:** `%t`\n"+
				user.Snowflake, username, user.FetchData,
		)
	}

	return output
}
