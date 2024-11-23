package guildData

import (
	"log"
	"os"
	"time"

	"github.com/bwmarrin/discordgo"
)

type Snowflake string

type GuildUser struct {
	Snowflake       Snowflake
	DiscordUsername string
	McUsername      string
	McUUID          string
}

type GuildBot struct {
	Users          map[Snowflake]GuildUser
	Events         map[string]Event
	EventSaver     EventSaver
	DiscordSession *discordgo.Session
}

func (bot *GuildBot) StartEventUpdater(ticker time.Ticker) {
	// Use a goroutine to check the events periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Checking for updates at %v", time.Now().UTC())
				bot.checkAndUpdateEvents()
			}
		}
	}()
}

func (bot *GuildBot) EndEventUpdater(ticker time.Ticker) {
	// Use a goroutine to check the events periodically
	go func() {
		for {
			select {
			case <-ticker.C:
				log.Printf("Checking for end event at %v", time.Now().UTC())
				bot.checkForEnd()
			}
		}
	}()
}

func (bot *GuildBot) checkForEnd() {
	if bot.EventSaver == nil {
		log.Printf("EventSaver is not set; cannot update events")
		return
	}
	for _, event := range bot.Events {

		if !event.GetIsActive() || event.HasEnded() {
			return
		}
		event.End()
	}
}

// checkAndUpdateEvents iterates through all events and performs updates if needed
func (bot *GuildBot) checkAndUpdateEvents() {
	if bot.EventSaver == nil {
		log.Printf("EventSaver is not set; cannot update events")
		return
	}

	// Iterate through all events and attempt to save the data
	for _, event := range bot.Events {

		if !event.GetIsActive() {
			return
		}

		// Update the data
		if event.ShouldFetchData() {
			log.Printf("Updating event (%s) with new data:", event.GetId())
			err := event.FetchData()
			if err != nil {
				log.Printf("Failed to fetch data: %v", err)
			}

			// Attempt to save the event data to the database
			err = bot.EventSaver.SaveEventData(event)
			if err != nil {
				log.Printf("Failed to save event data for event ID: %v, Type: %v. Error: %v", event.GetId(), event.GetType(), err)
			} else {
				log.Printf("Successfully saved event data for event ID: %v, Type: %v", event.GetId(), event.GetType())
			}
			event.SetLastFetch(time.Now().UTC())
		}

		if !event.IsHidden() {
			log.Printf("Showing event data for event: %s", event.GetId())
			// Send leaderboard update to the specified channel
			channelID := os.Getenv("DISCORD_CHANNEL_ID")
			leaderboardEmbed := GetLeaderboard(event)
			if leaderboardEmbed != nil {
				_, err := bot.DiscordSession.ChannelMessageSendEmbed(channelID, leaderboardEmbed)
				if err != nil {
					log.Printf("Failed to send leaderboard for event ID: %v, Error: %v", event.GetId(), err)
				} else {
					log.Printf("Successfully sent leaderboard for event ID: %v", event.GetId())
				}
			}
		}
	}
}
