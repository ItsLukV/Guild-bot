package models

type User struct {
	McUUID string
}

type GuildBot struct {
	Users map[string]User
}

func (g *GuildBot) Save() {
	// TODO
}
