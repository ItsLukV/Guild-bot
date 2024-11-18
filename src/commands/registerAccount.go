package commands

import (
	"fmt"
	"log"
	"strings"

	"github.com/ItsLukV/Guild-bot/src/db"
	guildData "github.com/ItsLukV/Guild-bot/src/guildData"
	"github.com/ItsLukV/Guild-bot/src/utils"
	"github.com/bwmarrin/discordgo"
)

func registerAccount(g *guildData.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate) {
	// Check if the user is already registered
	if v, exists := g.Users[i.Member.User.ID]; exists {
		utils.RespondWithError(s, i, fmt.Sprintf("You are already registered with the Minecraft IGN: %s", v.McUsername))
		log.Println("User already registered:", i.Member.User.ID)
		return
	}

	// Retrieve the Minecraft username from the command options
	userName := i.ApplicationCommandData().Options[0].StringValue()

	// Get the Minecraft UUID for the provided username
	mcUuid, err := utils.GetMCUUID(userName)
	if err != nil {
		log.Println("Error getting the Minecraft UUID:", err)
		var out string
		switch err.Error() {
		case "error: Unable to fetch data. Status code: 400":
			out = "Invalid Minecraft account name (Not possible)"
		case "error: Unable to fetch data. Status code: 404":
			out = "Invalid Minecraft account name (No existing account with that name)"
		default:
			out = "Unknown error, please contact staff"
		}

		utils.RespondWithError(s, i, out)
		return
	}

	// Check the Discord username associated with the Minecraft UUID
	discordName, err := utils.CheckUserName(mcUuid.Id)
	if err != nil {
		log.Println("Error fetching username:", err)
		utils.RespondWithError(s, i, "Failed to fetch username, please try again later.")
		return
	}

	// Verify that the Discord username matches
	if discordName != nil && strings.EqualFold(*discordName, i.Member.User.Username) {
		// Register the user
		g.Users[i.Member.User.ID] = guildData.GuildUser{
			Snowflake:       i.Member.User.ID,
			McUUID:          mcUuid.Id,
			McUsername:      mcUuid.Name,
			DiscordUsername: i.Member.User.Username,
		}
		user := g.Users[i.Member.User.ID]
		err := db.GetInstance().SaveUser(user)
		if err != nil {
			log.Println("Error saving user:", err)
			utils.RespondWithError(s, i, "Failed to save user (Please contact support)")
			return
		} else {
			// Registration successful, respond to the user
			err = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
				Type: discordgo.InteractionResponseChannelMessageWithSource,
				Data: &discordgo.InteractionResponseData{
					Content: "Registration complete!",
				},
			})
			if err != nil {
				log.Println("Error responding to interaction:", err)
			}
			return
		}
	} else {
		utils.RespondWithError(s, i, "You are not registered to this account")
		return
	}
}
