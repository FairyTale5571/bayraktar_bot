package models

import (
	"fmt"
	"time"
)

type Config struct {
	SteamKey      string `env:"STEAMAPI_KEY,required"`
	DiscordToken  string `env:"DISCORD_TOKEN,required"`
	TelegramToken string `env:"TELEGRAM_TOKEN"`

	VipRole   string `env:"VIP_ROLE_ID,required"`
	GuildID   string `env:"GUILD_ID,required"`
	RegRoleID string `env:"REG_ROLE_ID,required"`

	PostPassword string `env:"POST_PASSWORD,required"`

	MysqlUri string `env:"MYSQL_URI,required"`
	RedisUri string `env:"REDISCLOUD_URL,required"`

	URL   string `env:"URL,required"`
	PORT  string `env:"PORT" envDefault:"3000"`
	Debug bool   `env:"DEBUG,required"`
}

type News struct {
	ID          int
	Title       string
	Description string
	Link        string
	Published   time.Time
}

type NewsArray struct {
	News []News
}

func (n News) String() string {
	return fmt.Sprintf("[\"%s\", \"%s\", \"%s\", \"%s\"]", n.Title, n.Description, n.Link, n.Published.Format("15:04:05 01-02-2006"))
}

func (n NewsArray) MakeArmaArray() string {
	var res = "["
	for _, news := range n.News {
		res += news.String() + ","

	}
	res = res[:len(res)-1]
	res += "]"
	return res
}
