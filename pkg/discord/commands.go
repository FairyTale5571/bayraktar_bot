package discord

import "github.com/bwmarrin/discordgo"

var (
	commands = []*discordgo.ApplicationCommand{
		{
			Name:        "help",
			Description: "Бот для администрирования сервера Rocket, функционал доступен только администраторам",
		},
		{
			Name:        "help-player",
			Description: "Много ответов на много вопросов",
		},
		{
			Name:        "re-role",
			Description: "Перепроверяет выданные роли",
		},
		{
			Name:        "get-him",
			Description: "Получить данные игрока",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "Тегните пользователя",
					Required:    true,
				},
			},
		},
		{
			Name:        "give-boost",
			Description: "Выдать подарок за буст сервера",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionUser,
					Name:        "user-option",
					Description: "Тегните пользователя",
					Required:    true,
				},
			},
		},
		{
			Name:        "copy-role",
			Description: "Скопировать роль",
			Options: []*discordgo.ApplicationCommandOption{
				{
					Type:        discordgo.ApplicationCommandOptionRole,
					Name:        "role-option",
					Description: "Выберите роль",
					Required:    true,
				},
				{
					Type:        discordgo.ApplicationCommandOptionString,
					Name:        "role-name",
					Description: "Установите имя для новой роли",
					Required:    true,
				},
			},
		},
	}
)
