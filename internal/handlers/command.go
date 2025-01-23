package handlers

import (
	"github.com/ItsLukV/Guild-bot/internal/handlers/autocompletions"
	"github.com/ItsLukV/Guild-bot/internal/handlers/commands"
	"github.com/ItsLukV/Guild-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

var (
	Commands = []*discordgo.ApplicationCommand{
		{
			Name:        "fetch_users",
			Description: "Fetches all users",
		},
		{
			Name:        "fetch_guild_events",
			Description: "Fetches all guild events",
		},
		{
			Name:        "fetch_guild_event",
			Description: "Fetches a specific guild event",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "event_id",
					Description:  "The ID of the event to fetch",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
		{
			Name:        "fetch_user",
			Description: "Fetches a specific user and their event data",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:         discordgo.ApplicationCommandOptionString,
					Name:         "user_name",
					Description:  "The ID of the user to fetch",
					Required:     true,
					Autocomplete: true,
				},
			},
		},
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate, pm *utils.PaginatedSessions){
		"fetch_users":        commands.FetchUsersCommand,
		"fetch_guild_events": commands.FetchGuildEventsCommand,
		"fetch_guild_event":  commands.FetchGuildEventCommand,
		"fetch_user":         commands.FetchUserCommand,
	}

	AutocompleteHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fetch_guild_event": autocompletions.GuildEventAutocomplete,
		"fetch_user":        autocompletions.UserAutocomplete,
	}
)
