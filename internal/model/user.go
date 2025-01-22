package model

import (
	"encoding/json"
	"fmt"
	"log"
	"net/http"

	"github.com/ItsLukV/Guild-bot/internal/config"
)

type User struct {
	Id                string `json:"id"`
	ActiveProfileUUID string `json:"active_profile_UUID"`
	Snowflake         string `json:"discord_snowflake"`
	FetchData         bool   `json:"fetch_data"`
	IsInline          bool
}

func (u User) GetSectionName() string {
	return "**User data**"
}

func (u User) GetSectionValue() string {
	var output string
	username, err := FetchUsername(u.Id)

	if err != nil {
		log.Println("Error fetching username:", err)
		username = "Unknown"
	}

	output += fmt.Sprintf(
		"**Users name: <@%s>**\n"+
			"**Minecraft username:** `%s`\n"+
			"**Fetching Data:** `%t`\n\n",
		u.Snowflake, username, u.FetchData,
	)

	return output
}

func (u User) GetSectionInline() bool {
	return false
}

func (u *User) SetInLine(inline bool) {
	u.IsInline = inline
}

type UserWithEventData struct {
	DianaData    DianaData    `json:"diana_data"`
	DungeonsData DungeonsData `json:"dungeons_data"`
	User         User         `json:"user"`
}

type UsersResponse struct {
	Users []User `json:"users"`
}

func FetchUsername(uuid string) (string, error) {
	url := fmt.Sprintf("https://api.minecraftservices.com/minecraft/profile/lookup/%s", uuid)
	resp, err := http.Get(url)
	if err != nil {
		return "", err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return "", fmt.Errorf("failed to fetch username, status code: %d", resp.StatusCode)
	}

	var result struct {
		Name string `json:"name"`
	}
	if err := json.NewDecoder(resp.Body).Decode(&result); err != nil {
		return "", err
	}

	return result.Name, nil
}

func FetchUsers() ([]User, error) {
	url := fmt.Sprintf("%s/api/users", config.GlobalConfig.ApiBaseURL)
	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch users, status code: %d", resp.StatusCode)
	}

	var usersResponse UsersResponse
	if err := json.NewDecoder(resp.Body).Decode(&usersResponse); err != nil {
		return nil, err
	}

	return usersResponse.Users, nil
}

func FetchUser(userID string) (*UserWithEventData, error) {
	url := fmt.Sprintf("%s/api/user?id=%s", config.GlobalConfig.ApiBaseURL, userID)

	resp, err := http.Get(url)
	if err != nil {
		return nil, err
	}
	defer resp.Body.Close()

	if resp.StatusCode != http.StatusOK {
		return nil, fmt.Errorf("failed to fetch user data, status code: %d", resp.StatusCode)
	}

	var fud UserWithEventData
	if err := json.NewDecoder(resp.Body).Decode(&fud); err != nil {
		return nil, err
	}

	return &fud, nil
}
