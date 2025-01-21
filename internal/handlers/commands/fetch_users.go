package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/ItsLukV/Guild-bot/internal/config"
	models "github.com/ItsLukV/Guild-bot/internal/model"
	"github.com/ItsLukV/Guild-bot/internal/restclient"
	"github.com/bwmarrin/discordgo"
)

func FetchUsersCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url := fmt.Sprintf("%s/api/users", config.GlobalConfig.ApiBaseURL)
	log.Println("Fetching data from:", url)
	data, err := restclient.FetchApi[models.UsersResponse](url)
	if err != nil {
		log.Println("Error fetching data:", err)
		s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Failed to fetch data.",
			},
		})
		return
	}

	var sb strings.Builder
	sb.WriteString("Fetched users:\n")
	for _, user := range data.Users {
		sb.WriteString(fmt.Sprintf("ID: %s, Active Profile UUID: %s, Fetch Data: %t\n", user.ID, user.ActiveProfileUUID, user.FetchData))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}
