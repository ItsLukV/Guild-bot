package commands

import (
	"github.com/ItsLukV/Guild-bot/src/db"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name: "register",
			// All commands and options must have a description
			// Commands/options without description will fail the registration
			// of the command.
			Description: "Link your minecraft account with your discord account",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "minecraft-username",
					Description: "Minecraft username to be added",
					Required:    true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(g *db.GuildBot, s *discordgo.Session, i *discordgo.InteractionCreate){
		"register": registerAccount,
	}
)
