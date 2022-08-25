package discord

import (
	"fmt"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) components() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"login":         d.componentLogin,
		"how_to_play":   d.howToPlay,
		"create_ticket": d.createTicket,
		"close_ticket":  d.closeTicket,
	}
}

func (d *Discord) createTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: "думаю...",
		},
	})
	if err != nil {
		d.logger.Errorf("createTicket(): Error interaction respond: %s", err.Error())
		return
	}
	ch, err := d.createTickets(i.GuildID, i.Interaction.Member.User)
	if err != nil {
		d.logger.Errorf("createTicket(): Error create tickets: %s", err.Error())
		return
	}
	resText := fmt.Sprintf("Тикет был создан в канале <#%s>", ch.ID)
	_, err = d.ds.InteractionResponseEdit(i.Interaction, &discordgo.WebhookEdit{
		Content: &resText,
	})
	if err != nil {
		d.logger.Errorf("createTicket(): Error interaction response edit: %s", err.Error())
		return
	}
}

func (d *Discord) printLogin(id string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Авторизация",
		Description: "Для синхронизации сервера с дискордом, вам необходимо:",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Зайти на сервер", Value: "С подробной инструкцией можно ознакомится в канале <#872230767625908285>"},
			{Name: "Нажать на кнопку снизу", Value: "И авторизоваться через Steam по ссылке в личном сообщении"},
		},
		Color: 0x8700ff,
	}
	data := &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Залогиниться",
						Style:    discordgo.SuccessButton,
						Disabled: false,
						CustomID: "login",
					},
				},
			},
		},
	}

	_, err := d.ds.ChannelMessageSendComplex(id, data)
	if err != nil {
		d.logger.Errorf("printLogin(): Error sending message: %s", err.Error())
		return
	}
}

func (d *Discord) componentLogin(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := s.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: "Дальнейшие инструкции отправлена вам в личные сообщения!",
		},
	})
	embed := &discordgo.MessageEmbed{
		Title:       "Авторизация",
		Description: "Для авторизации нажми на кнопку \"Верифицировать\"",
		Color:       0x8700ff,
	}
	data := &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Верифицировать",
						Style:    discordgo.LinkButton,
						Disabled: false,
						URL:      d.steam.GetAuthLink(i.GuildID, i.Interaction.Member.User.ID),
					},
				},
			},
		},
	}

	ch, err := d.ds.UserChannelCreate(i.Interaction.Member.User.ID)
	if err != nil {
		d.logger.Errorf("componentLogin(): Error user channel create: %s", err.Error())
		return
	}
	_, err = d.ds.ChannelMessageSendComplex(ch.ID, data)
	if err != nil {
		d.logger.Errorf("componentLogin(): Error sending message: %s", err.Error())
		return
	}
}

func (d *Discord) howToPlay(s *discordgo.Session, i *discordgo.InteractionCreate) {
	embed, component := d.getHow2Play()
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: component,
			Embeds: []*discordgo.MessageEmbed{
				embed,
			},
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		d.logger.Errorf("howToPlay(): Error responding: %s", err.Error())
	}
}
