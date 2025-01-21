package discord

import (
	"log"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/handlers"
	"github.com/bwmarrin/discordgo"
)

type Service struct {
	session            *discordgo.Session
	registeredCommands []*discordgo.ApplicationCommand
}

func New(token string) *Service {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	service := &Service{
		session: session,
	}

	return service
}

func (s *Service) AddCommandHandlers() {
	s.session.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {

		case discordgo.InteractionApplicationCommand:
			if h, ok := handlers.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(sess, i)
			}

		case discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := handlers.AutocompleteHandlers[i.ApplicationCommandData().Name]; ok {
				h(sess, i)
			}

		case discordgo.InteractionMessageComponent:
			handlePaginationButton(sess, i)

		default:
			log.Panicf("unexpected discordgo.InteractionType: %#v", i.Interaction.Type)
		}
	})
}

func (s *Service) Start() error {
	return s.session.Open()
}

func (s *Service) RegisterCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, len(handlers.Commands))
	for i, v := range handlers.Commands {
		var c *discordgo.ApplicationCommand
		var err error
		if config.GlobalConfig.GuildID != "" {
			c, err = s.session.ApplicationCommandCreate(
				s.session.State.User.ID,
				config.GlobalConfig.GuildID,
				v,
			)
		} else {
			c, err = s.session.ApplicationCommandCreate(
				s.session.State.User.ID,
				"",
				v,
			)
		}
		if err != nil {
			log.Panicf("Cannot create '%v' command: %v", v.Name, err)
		}
		cmds[i] = c
	}
	s.registeredCommands = cmds
	return cmds
}

func (s *Service) RemoveCommands(cmds []*discordgo.ApplicationCommand) {
	for _, v := range cmds {
		var err error
		if config.GlobalConfig.GuildID != "" {
			err = s.session.ApplicationCommandDelete(
				s.session.State.User.ID,
				config.GlobalConfig.GuildID,
				v.ID,
			)
		} else {
			err = s.session.ApplicationCommandDelete(
				s.session.State.User.ID,
				"",
				v.ID,
			)
		}
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (s *Service) Session() *discordgo.Session {
	return s.session
}
