package discord

import (
    "log"
    "strings"

    "github.com/ItsLukV/Guild-bot/internal/utils"
    "github.com/bwmarrin/discordgo"
)

func handlePaginationButton(s *discordgo.Session, i *discordgo.InteractionCreate) {
    data := i.MessageComponentData()
    customID := data.CustomID // e.g. "page_prev_<uuid>" or "page_next_<uuid>"

    parts := strings.SplitN(customID, "_", 3)
    if len(parts) < 3 {
        // not a pagination ID
        return
    }
    direction := parts[1] // "prev" or "next"
    paginationID := parts[2]

    // Look up pagination data
    pd, ok := utils.PaginationStore[paginationID]
    if !ok {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "Pagination data not found or expired.",
            },
        })
        return
    }

    // (Optional) If you only want the original author to page through, check:
    if i.Member.User.ID != pd.AuthorID {
        s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
            Type: discordgo.InteractionResponseChannelMessageWithSource,
            Data: &discordgo.InteractionResponseData{
                Content: "You are not allowed to change pages.",
                Flags:   discordgo.MessageFlagsEphemeral,
            },
        })
        return
    }

    // Update page
    if direction == "prev" && pd.PageIndex > 0 {
        pd.PageIndex--
    } else if direction == "next" && pd.PageIndex < len(pd.Pages)-1 {
        pd.PageIndex++
    }

    // Rebuild the embed
    embed := utils.MakePaginationEmbed(pd)

    // Rebuild the components
    comps := utils.MakePaginationComponents(paginationID, pd.PageIndex, len(pd.Pages))

    // Use InteractionResponseUpdateMessage to edit in place
    err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
        Type: discordgo.InteractionResponseUpdateMessage,
        Data: &discordgo.InteractionResponseData{
            Embeds:     []*discordgo.MessageEmbed{embed},
            Components: comps,
        },
    })
    if err != nil {
        log.Println("Error responding to pagination button:", err)
    }
}
