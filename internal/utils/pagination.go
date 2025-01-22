package utils

import (
	"fmt"
	"log"
	"math"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// PaginationData holds pages and current state for a single pagination session.
type PaginationData struct {
	Fields      []Section
	PageSize    int
	PageIndex   int
	Description string
	AuthorID    string
	Title       string
	Footer      string
	Color       int
	CreatedAt   time.Time
}

func (p *PaginationData) GetPageAmount() int {
	return int(math.Ceil(float64(len(p.Fields)) / float64(p.PageSize)))
}

// We store them in a global map: paginationID -> *PaginationData
var PaginationStore = map[string]*PaginationData{}

// BuildPaginationID creates a unique ID for this pagination session
func BuildPaginationID() string {
	return uuid.NewString()
}

// MakePaginationComponents builds a row of "prev" / "next" buttons
func MakePaginationComponents(paginationID string, pageIndex, totalPages int) []discordgo.MessageComponent {
	prevDisabled := pageIndex <= 0
	nextDisabled := pageIndex >= totalPages-1

	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label:    "Previous",
					Style:    discordgo.PrimaryButton,
					CustomID: "page_prev_" + paginationID,
					Disabled: prevDisabled,
				},
				discordgo.Button{
					Label:    "Next",
					Style:    discordgo.PrimaryButton,
					CustomID: "page_next_" + paginationID,
					Disabled: nextDisabled,
				},
			},
		},
	}
}

// MakePaginationEmbed builds a simple embed for the specified page
func MakePaginationEmbed(data *PaginationData) *discordgo.MessageEmbed {
	var fields []*discordgo.MessageEmbedField

	for i := 0; i < data.PageSize; i++ {
		// If there is no more fields to add break
		if len(data.Fields)-1 < data.PageIndex*data.PageSize+i {
			break
		}

		// Append the field
		field := data.Fields[data.PageIndex*data.PageSize+i]
		fields = append(fields, &discordgo.MessageEmbedField{
			Name:   field.GetSectionName(),
			Value:  field.GetSectionValue(),
			Inline: field.GetSectionInline(),
		})
	}

	return &discordgo.MessageEmbed{
		Title:       data.Title,
		Description: data.Description,
		Color:       data.Color,
		Footer: &discordgo.MessageEmbedFooter{
			Text: fmt.Sprintf("Page: %v/%v", data.PageIndex, math.Ceil(float64(len(data.Fields))/5.0)),
		},
		Timestamp: time.Now().Format(time.RFC3339),
		Fields:    fields,
	}
}

func SendInitialPaginationResponse(
	s *discordgo.Session,
	i *discordgo.InteractionCreate,
	paginationID string,
	data *PaginationData,
) error {
	embed := MakePaginationEmbed(data)

	totalPages := data.GetPageAmount()
	if totalPages < 1 {
		totalPages = 1 // At least 1 page even if no fields
	}

	components := MakePaginationComponents(paginationID, data.PageIndex, totalPages)

	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Embeds:     []*discordgo.MessageEmbed{embed},
			Components: components,
		},
	})
	if err != nil {
		log.Println("Failed to respond with paginated embed:", err)
	}

	return err
}
