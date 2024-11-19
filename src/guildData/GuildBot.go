package guildData

type Snowflake string

type GuildUser struct {
	Snowflake       Snowflake
	DiscordUsername string
	McUsername      string
	McUUID          string
}

type GuildBot struct {
	Users  map[Snowflake]GuildUser
	Events map[string]Event
}
