package model

import (
	"encoding/json"
	"fmt"
	"github.com/ItsLukV/Guild-bot/internal/config"
	"github.com/ItsLukV/Guild-bot/internal/restclient"
	"net/http"
)

func FetchUsers() ([]User, error) {
	url := fmt.Sprintf("%s/api/users", config.GlobalConfig.ApiBaseURL)
	// Reuse generic fetch
	data, err := restclient.FetchApi[struct {
		Users []User `json:"users"`
	}](url)
	if err != nil {
		return nil, err
	}
	return data.Users, nil
}

func FetchUser(userID string) (*UserWithEventData, error) {
	url := fmt.Sprintf("%s/api/user?id=%s", config.GlobalConfig.ApiBaseURL, userID)
	return restclient.FetchApi[UserWithEventData](url)
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

func FetchUsername(uuid string) (string, error) {
	url := fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/lookup/%s", uuid)
	data, err := restclient.FetchApi[struct {
		Name string `json:"name"`
	}](url)
	if err != nil {
		return "", err
	}
	return data.Name, nil
}
