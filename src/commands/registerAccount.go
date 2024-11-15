package commands

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ItsLukV/Guild-bot/src/db"
	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func registerAccount(g *db.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	if v, exists := g.Users[i.Member.User.ID]; exists {
		err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("You are already registered, with the MC IGN: %s", v.McUsername),
			},
		})
		if err != nil {
			log.Println("Error responding to interaction:", err)
			return
		}
		log.Println("User already registered: ", i.Member.User.ID)
		return
	}

	// Access options in the order provided by the user.
	// There is only one option, so we just set it as the userName
	userName := i.ApplicationCommandData().Options[0].StringValue()

	mcUuid, err := utils.GetMCUUID(userName)
	if err != nil {
		log.Println("Error getting the minecraft UUID: ", err)
		var out string
		switch err.Error() {
		case "error: Unable to fetch data. Status code: 400":
			out = "Invalid minecraft account name (Not possible)"
		case "error: Unable to fetch data. Status code: 404":
			out = "Invalid minecraft account name (No exiting account with that name)"
		default:
			out = "Unknown error, please contact staff"
		}

		fmt.Println(reflect.TypeOf(err))
		err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: out,
			},
		})
		return
	}

	dicordName, err := utils.CheckUserName(mcUuid.Id)
	if err != nil {
		log.Println("Error fetching username: ", err)
		return
	}

	var out = "You are not registered to this account"
	if &i.Member.User.Username != nil &&
		strings.ToLower(*dicordName) == strings.ToLower(i.Member.User.Username) {
		g.Users[i.Member.User.ID] = db.GuildUser{
			Snowflake:       i.Member.User.ID,
			McUUID:          mcUuid.Id,
			McUsername:      mcUuid.Name,
			DiscordUsername: i.Member.User.Username,
		}
		user := g.Users[i.Member.User.ID]
		err := user.SaveUser()
		if err != nil {
			log.Println("Error saving user: ", err)
			out = "Failed to save user (Please contact support)"
		} else {
			out = "Registration complete"
		}
	}

	err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: out,
		},
	})

	if err != nil {
		log.Println("Error responding to interaction:", err)
		return
	}
}
