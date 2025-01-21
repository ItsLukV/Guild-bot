package discord

import (
	"log"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/handlers"
	"github.com/bwmarrin/discordgo"
)

type Service struct {
	session            *discordgo.Session
	apiBaseURL         string
	apiKey             string
	registeredCommands []*discordgo.ApplicationCommand
}

func New(token, apiBaseURL, apiKey string) *Service {
	session, err := discordgo.New("Bot " + token)
	if err != nil {
		log.Fatal("error creating Discord session,", err)
	}

	service := &Service{
		session:    session,
		apiBaseURL: apiBaseURL,
		apiKey:     apiKey,
	}

	session.AddHandler(func(s *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {
		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if handler, ok := handlers.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				handler(s, i)
			}
		}
	})

	return service
}

func (s *Service) Start() error {
	return s.session.Open()
}

func (s *Service) RegisterCommands() []*discordgo.ApplicationCommand {
	cmds := make([]*discordgo.ApplicationCommand, len(handlers.Commands))
	for i, v := range handlers.Commands {
		c, err := s.session.ApplicationCommandCreate(
			s.session.State.User.ID,
			config.GlobalConfig.GuildID,
			v,
		)
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
		err := s.session.ApplicationCommandDelete(
			s.session.State.User.ID,
			config.GlobalConfig.GuildID,
			v.ID,
		)
		if err != nil {
			log.Panicf("Cannot delete '%v' command: %v", v.Name, err)
		}
	}
}

func (s *Service) AddCommandHandlers() {
	s.session.AddHandler(func(sess *discordgo.Session, i *discordgo.InteractionCreate) {
		switch i.Type {

		case discordgo.InteractionApplicationCommand, discordgo.InteractionApplicationCommandAutocomplete:
			if h, ok := handlers.CommandHandlers[i.ApplicationCommandData().Name]; ok {
				h(sess, i)
			}
		default:
			log.Panicf("unexpected discordgo.InteractionType: %#v", i.Interaction.Type)
		}
	})
}
