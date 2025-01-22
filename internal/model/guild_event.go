package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"
	"time"

	"github.com/ItsLukV/Guild-bot/internal/config"
)

type GuildEvent struct {
	Id        string         `json:"id"`
	Users     []string       `json:"users"`
	StartTime time.Time      `json:"start_time"`
	Duration  int            `json:"duration"`
	Type      GuildEventType `json:"type"`
	IsHidden  bool           `json:"is_hidden"`
	EventData []EventData    `json:"event_data"`
	IsInLine  bool
}

type GuildEventType string

const (
	Dungeons GuildEventType = "dungeons"
	Diana    GuildEventType = "diana"
)

func (ge *GuildEvent) GetUserNames() ([]string, error) {
	var usernames []string
	for _, uuid := range ge.Users {
		username, err := FetchUsername(uuid)
		if err != nil {
			return nil, err
		}
		usernames = append(usernames, username)
	}
	return usernames, nil
}

func (g *GuildEvent) UnmarshalJSON(data []byte) error {
	// Create a local struct alias that matches GuildEvent but uses
	// json.RawMessage for event_data so we can parse it manually.
	type rawGuildEvent struct {
		ID        string          `json:"id"`
		Users     []string        `json:"users"`
		StartTime time.Time       `json:"start_time"`
		Duration  int             `json:"duration"`
		Type      GuildEventType  `json:"type"`
		IsHidden  bool            `json:"is_hidden"`
		EventData json.RawMessage `json:"event_data"` // We'll decode this manually
	}

	var raw rawGuildEvent
	if err := json.Unmarshal(data, &raw); err != nil {
		return err
	}

	// Copy over basic fields
	g.Id = raw.ID
	g.Users = raw.Users
	g.StartTime = raw.StartTime
	g.Duration = raw.Duration
	g.Type = raw.Type
	g.IsHidden = raw.IsHidden

	if len(raw.EventData) == 0 || string(raw.EventData) == "null" {
		// No event_data, skip decoding
		g.EventData = nil
		return nil
	}

	// Switch on the event type to decode the event_data
	switch g.Type {
	case Diana:
		// Decode as []DianaData
		var dianaSlice []DianaData
		if err := json.Unmarshal(raw.EventData, &dianaSlice); err != nil {
			return err
		}

		// Convert []DianaData to []EventData
		g.EventData = make([]EventData, len(dianaSlice))
		for i := range dianaSlice {
			g.EventData[i] = &dianaSlice[i]
		}

	case Dungeons:
		// Decode as []DungeonsData
		var dungeonSlice []DungeonsData
		if err := json.Unmarshal(raw.EventData, &dungeonSlice); err != nil {
			return err
		}

		// Convert []DungeonsData to []EventData
		g.EventData = make([]EventData, len(dungeonSlice))
		for i := range dungeonSlice {
			g.EventData[i] = &dungeonSlice[i]
		}

	default:
		return fmt.Errorf("unknown event type: %s", g.Type)
	}

	return nil
}

func (ge *GuildEvent) GetSectionName() string {
	return "**Guild Event**"
}

func (ge *GuildEvent) GetSectionValue() string {
	var output string

	// Convert user IDs to names
	userNames, err := ge.GetUserNames()
	if err != nil {
		log.Printf("Error getting user names: %v\n", err)
		userNames = ge.Users
	}

	eventText := fmt.Sprintf(
		"Event ID: `%s`\n"+
			"Users: `%v`\n"+
			"Start Time: `%s`\n"+
			"Duration: `%dh`\n"+
			"Type: `%s`\n"+
			"Hidden: `%t`\n\n",
		ge.Id,
		userNames,
		ge.StartTime.Format(time.RFC1123),
		ge.Duration,
		ge.Type,
		ge.IsHidden,
	)

	output += eventText

	return output
}

func (ge *GuildEvent) GetSectionInline() bool {
	return ge.IsInLine
}

func (ge *GuildEvent) SetInLine(IsInLine bool) {
	ge.IsInLine = IsInLine
}

func FetchGuildEvent(eventId string) (*GuildEvent, error) {
	url := fmt.Sprintf("%s/api/guildevent?id=%s", config.GlobalConfig.ApiBaseURL, eventId)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch guild event, status code: %d", resp.StatusCode)
	}

	var guildEvent GuildEvent
	if err := json.NewDecoder(resp.Body).Decode(&guildEvent); err != nil {
		return nil, err
	}

	return &guildEvent, nil
}

func FetchGuildEvents() ([]GuildEvent, error) {
	url := fmt.Sprintf("%s/api/guildevents", config.GlobalConfig.ApiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch guild events, status code: %d", resp.StatusCode)
	}

	var guildEvents []GuildEvent
	if err := json.NewDecoder(resp.Body).Decode(&guildEvents); err != nil {
		return nil, err
	}

	return guildEvents, nil
}
