package guildData

type GuildUser struct {
	Snowflake       string
	McUUID          string
	McUsername      string
	DiscordUsername string
}

type GuildBot struct {
	Users  map[string]GuildUser
	Events map[int]Event
}
