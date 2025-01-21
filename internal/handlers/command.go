package handlers

import (
	"github.com/ItsLukV/Guild-bot/internal/handlers/autocompletions"
	"github.com/ItsLukV/Guild-bot/internal/handlers/commands"
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
	}

	CommandHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"fetch_users":        commands.FetchUsersCommand,
		"fetch_guild_events": commands.FetchGuildEventsCommand,
		"fetch_guild_event":  commands.FetchGuildEventCommand,
	}

	AutocompleteHandlers = map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		"fetch_guild_event": autocompletions.GuildEventAutocomplete,
	}
)
