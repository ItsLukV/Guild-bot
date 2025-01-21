package model

import (
	"encoding/json"
	"fmt"
	"net/http"
)

type User struct {
	ID                string `json:"id"`
	ActiveProfileUUID string `json:"active_profile_UUID"`
	FetchData         bool   `json:"fetch_data"`
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
