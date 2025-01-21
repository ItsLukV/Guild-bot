package commands

import (
	"fmt"
	"log"
	"time"

	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/model"
	"github.com/ItsLukV/Guild-bot/internal/utils"
	"github.com/bwmarrin/discordgo"
)

func FetchGuildEventCommand(s *discordgo.Session, i *discordgo.InteractionCreate) {
	eventID := i.ApplicationCommandData().Options[0].StringValue()

	guildEvent, err := model.FetchGuildEvent(eventID)
	if err != nil {
		log.Printf("Error fetching event: %v", err)
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: "Error fetching event data.",
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	if guildEvent == nil || guildEvent.Id == "" {
		_ = s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
			Type: discordgo.InteractionResponseChannelMessageWithSource,
			Data: &discordgo.InteractionResponseData{
				Content: fmt.Sprintf("No event found for ID: %s", eventID),
				Flags:   discordgo.MessageFlagsEphemeral,
			},
		})
		return
	}

	// Build a large string for the event
	fullText := buildLargeEventString(guildEvent)

	// Chunk into ~1000-char pages
	pages := utils.ChunkString(fullText, 1000)
	if len(pages) == 0 {
		pages = []string{"No data available."}
	}

	// Create a unique pagination ID
	paginationID := utils.BuildPaginationID()

	// Create and store the PaginationData
	utils.PaginationStore[paginationID] = &utils.PaginationData{
		Pages:     pages,
		PageIndex: 0,
		AuthorID:  i.Member.User.ID,
		Title:     "Fetched Guild Event",
		Footer:    fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		CreatedAt: time.Now(),
	}

	// Make the first embed
	embed := utils.MakePaginationEmbed(utils.PaginationStore[paginationID])

	// Buttons
	components := utils.MakePaginationComponents(paginationID, 0, len(pages))

	// Respond
	if err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	}); err != nil {
		log.Println("Failed to respond with guild event embed:", err)
	}
}

// Build a large string from the event data (like your "fieldValue" but unconstrained).
func buildLargeEventString(ev *model.GuildEvent) string {
	userNames, err := ev.GetUserNames()
	if err != nil {
		log.Printf("Error fetching user names: %v\n", err)
		userNames = ev.Users
	}

	baseInfo := fmt.Sprintf(
		"Event ID: %s\nUsers: %v\nStart: %s\nDuration: %dh\nType: %s\nHidden: %t\n\n",
		ev.Id,
		userNames,
		ev.StartTime.Format(time.RFC1123),
		ev.Duration,
		ev.Type,
		ev.IsHidden,
	)

	// If there's event data, append it
	if len(ev.EventData) == 0 {
		baseInfo += "No event data.\n"
		return baseInfo
	}

	baseInfo += "=== Event Data ===\n"
	for _, dataItem := range ev.EventData {
		switch ev.Type {
		case model.Diana:
			dianaData, ok := dataItem.(model.DianaData)
			if !ok {
				continue
			}
			// Convert ID to name
			userName, uErr := model.FetchUsername(dianaData.Id)
			if uErr != nil {
				userName = dianaData.Id
			}
			baseInfo += fmt.Sprintf(
				"- %s (Diana)\n  FetchTime: %s\n  BurrowsTreasure: %d\n  BurrowsCombat: %d\n  GaiaConstruct: %d\n  MinosChampion: %d\n  MinosHunter: %d\n  MinosInquisitor: %d\n  Minotaur: %d\n  SiameseLynx: %d\n\n",
				userName,
				dianaData.FetchTime.Format(time.RFC1123),
				dianaData.BurrowsTreasure,
				dianaData.BurrowsCombat,
				dianaData.GaiaConstruct,
				dianaData.MinosChampion,
				dianaData.MinosHunter,
				dianaData.MinosInquisitor,
				dianaData.Minotaur,
				dianaData.SiameseLynx,
			)
		case model.Dungeons:
			dungeonData, ok := dataItem.(model.DungeonsData)
			if !ok {
				continue
			}
			userName, uErr := model.FetchUsername(dungeonData.ID)
			if uErr != nil {
				userName = dungeonData.ID
			}
			baseInfo += fmt.Sprintf(
				"- %s (Dungeons)\n  FetchTime: %s\n  Experience: %.2f\n  Secrets: %d\n  Completions: %v\n  MasterCompletions: %v\n  ClassXP: %v\n\n",
				userName,
				dungeonData.FetchTime.Format(time.RFC1123),
				dungeonData.Experience,
				dungeonData.Secrets,
				dungeonData.Completions,
				dungeonData.MasterCompletions,
				dungeonData.ClassXP,
			)
		default:
			baseInfo += fmt.Sprintf("- Unknown Type: %v\n\n", dataItem)
		}
	}
	return baseInfo
}

// Convert a page of text into an embed
func makeEmbedPage(ev *model.GuildEvent, pageText string, pageIndex, totalPages int) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       "Fetched Guild Event",
		Description: fmt.Sprintf("Page %d / %d", pageIndex+1, totalPages),
		Color:       0x1F8B4C,
		Fields: []*discordgo.MessageEmbedField{
			{
				Name:  fmt.Sprintf("Event: %s", ev.Id),
				Value: pageText,
			},
		},
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Fetched from %s", config.GlobalConfig.ApiBaseURL),
		},
		Timestamp: time.Now().Format(time.RFC3339),
	}
}

// Build "Prev" / "Next" buttons.
func makePaginationComponents(paginationID string, pageIndex, totalPages int) []discordgo.MessageComponent {
	// If on first page, disable "prev" button
	prevDisabled := pageIndex <= 0
	nextDisabled := pageIndex >= totalPages-1

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Previous",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("page_prev_%s", paginationID),
					Disabled: prevDisabled,
				},
				discordgo.Button{
					Label:    "Next",
					Style:    discordgo.PrimaryButton,
					CustomID: fmt.Sprintf("page_next_%s", paginationID),
					Disabled: nextDisabled,
				},
			},
		},
	}
}
