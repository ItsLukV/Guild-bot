package db

import "log"

type GuildUser struct {
	Snowflake       string
	McUUID          string
	McUsername      string
	DiscordUsername string
}

func (g *GuildUser) SaveUser() error {
	err := GetInstance().AddUser(*g)
	if err != nil {
		log.Printf("Error adding user: %v", err)
		return err
	}
	return nil
}

type GuildBot struct {
	Users map[string]GuildUser
}

func (g *GuildBot) Save() {
	GetInstance().Save(*g)
}
