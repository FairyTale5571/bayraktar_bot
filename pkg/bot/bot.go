package bot

import (
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	tgbotapi "github.com/go-telegram-bot-api/telegram-bot-api/v5"
)

type Bot struct {
	logger *logger.LoggerWrapper

	bot     *tgbotapi.BotAPI
	updates tgbotapi.UpdatesChannel
}

func New() {

}

const (
	updateOffset  = 0
	updateTimeout = 64
)

func (b *Bot) Start() {

}
