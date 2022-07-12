package discord

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Бот для администрирования сервера Rimas, функционал доступен только администраторам",
		},
	}
)
