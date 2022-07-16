package models

import (
	"os"

	discordoauth "github.com/ravener/discord-oauth2"
	"golang.org/x/oauth2"
)

type Config struct {
	SteamKey      string `env:"STEAMAPI_KEY,required"`
	DiscordToken  string `env:"DISCORD_TOKEN,required"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`

	MysqlUri string `env:"MYSQL_URI,required"`
	RedisUri string `env:"REDISCLOUD_URL,required"`

	URL   string `env:"URL,required"`
	PORT  string `env:"PORT,required"`
	Debug bool   `env:"DEBUG,required"`
}

var (
	DiscordOauth = oauth2.Config{
		RedirectURL:  os.Getenv("URL") + "/auth/discord/callback",
		ClientID:     os.Getenv("DISCORD_CLIENT"),
		ClientSecret: os.Getenv("DISCORD_SECRET"),
		Scopes:       []string{discordoauth.ScopeIdentify, discordoauth.ScopeGuilds},
		Endpoint:     discordoauth.Endpoint,
	}
)
