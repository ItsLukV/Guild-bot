package model

import (
	"fmt"
	"log"
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
