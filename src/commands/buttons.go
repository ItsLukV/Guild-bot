package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func HandleButtonInteraction(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	str := strings.Split(i.MessageComponentData().CustomID, "_")
	switch str[0] {
	case "joinEvent":
		eventId := str[1]
		userId := i.Member.User.ID
		event, exists := g.Events[eventId]
		if !exists {
			utils.RespondWithErrorPrivate(s, i, "Event not found.")
			return
		}

		user, exists := g.Users[guildData.Snowflake(userId)]
		if !exists {
			utils.RespondWithErrorPrivate(s, i, "Please register your account with /register")
			return
		}
		_, ok := event.GetUsers()[user.Snowflake]
		if !ok {
			event.AddUser(&user)
		} else {
			utils.RespondWithErrorPrivate(s, i, "You are already resisted to the event.")
			return
		}

		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You have joined the event: **%s**!", event.GetType()),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		if err != nil {
			log.Printf("Failed to respond to button interaction: %v", err)
		}
	}
}
