package utils

import (
	"fmt"
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/google/uuid"
)

// chunkString: Utility to split a string into slices of maxLen
func ChunkString(s string, maxLen int) []string {
    var chunks []string
    for len(s) > maxLen {
        chunks = append(chunks, s[:maxLen])
        s = s[maxLen:]
    }
    if len(s) > 0 {
        chunks = append(chunks, s)
    }
    return chunks
}

// PaginationData holds pages and current state for a single pagination session.
type PaginationData struct {
    Pages     []string
    PageIndex int
    AuthorID  string // optional: track who can click
    Title     string // optional: some display data
    Footer    string // optional: some display data
    CreatedAt time.Time
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
    currentPage := data.Pages[data.PageIndex]
    totalPages := len(data.Pages)

    return &discordgo.MessageEmbed{
        Title:       data.Title,
        Description: currentPage,
        Color:       0x1F8B4C,
        Footer: &discordgo.MessageEmbedFooter{
            Text: data.Footer,
        },
        Timestamp: time.Now().Format(time.RFC3339),
        Fields: []*discordgo.MessageEmbedField{
            {
                Name:  "Page",
                Value: strings.TrimSpace(
                    // e.g. "Page 2 / 5"
                    PageIndicator(data.PageIndex, totalPages),
                ),
                Inline: false,
            },
        },
    }
}

// PageIndicator helper
func PageIndicator(pageIndex, total int) string {
    return  fmt.Sprintf("%d / %d", pageIndex+1, total)
}
