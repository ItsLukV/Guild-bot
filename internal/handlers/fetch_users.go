package handlers

import (
	"fmt"
	"log"
	"strings"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/restclient"
	"github.com/bwmarrin/discordgo"
)

type User struct {
	ID                 string `json:"id"`
	ActiveProfileUUID  string `json:"active_profile_UUID"`
	FetchData          bool   `json:"FetchData"`
}

func FetchUsersCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	url := fmt.Sprintf("%s/api/users", config.GlobalConfig.ApiBaseURL)
	data, err := restclient.FetchApi[[]User](url)
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
	for _, user := range *data {
		sb.WriteString(fmt.Sprintf("ID: %s, Active Profile UUID: %s, Fetch Data: %t\n", user.ID, user.ActiveProfileUUID, user.FetchData))
	}

	s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: sb.String(),
		},
	})
}
