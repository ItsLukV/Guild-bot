package handlers

import (
	"github.com/bwmarrin/discordgo"
	"github.com/ItsLukV/Guild-bot/internal/handlers/commands"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "fetch_users",
			Description: "Fetches all users",
		},
		{
			Name: 	  	 "fetch_guil_events",
			Description: "Fetches all guild events",
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fetch_users": commands.FetchUsersCommand,
		"fetch_guild_events": commands.FetchGuildEventsCommand,
	}
)
