package config

import (
	"log"
	"os"
)

type Config struct {
	BotToken   string
	ApiBaseURL string
	ApiKey     string
	GuildID    string
}

var GlobalConfig *Config

func LoadConfig() *Config {
	GlobalConfig = &Config{
		BotToken:   getEnv("BOT_TOKEN", ""),
		ApiBaseURL: getEnv("API_BASE_URL", ""),
		ApiKey:     getEnv("API_KEY", ""),
		GuildID:    getEnv("GUILD_ID", ""),
	}
	return GlobalConfig
}

func getEnv(key, fallback string) string {
	value := os.Getenv(key)
	if value != "" {
		return value
	}
	if fallback == "" {
		log.Fatalf("Environment variable %s not set and no fallback provided", key)
	}
	return fallback
}
