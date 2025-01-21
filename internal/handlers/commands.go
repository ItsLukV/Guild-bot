package handlers

import (
	"github.com/ItsLukV/Guild-bot/internal/handlers/commands"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "fetch_users",
			Description: "Fetches all users",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fetch_users": commands.FetchUsersCommand,
	}
)
