package commands

import (
	"fmt"
	"log"
	"reflect"
	"strings"

	"github.com/ItsLukV/Guild-bot/src/models"
	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func registerAccount(g *models.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
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
		g.Users[i.Member.User.Username] = models.User{McUUID: mcUuid.Id}
		g.Save()
		out = "Registration complete"
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
