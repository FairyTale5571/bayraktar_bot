package discord

import "github.com/bwmarrin/discordgo"

func (d *Discord) components() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"login": d.componentLogin,
	}
}

func (d *Discord) printLogin(id string) {
	data := &discordgo.MessageSend{
		Content: "Привет! Нажми на кнопку, чтобы залогиниться на сервере!",
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
	data := &discordgo.MessageSend{
		Content: "Привет! Авторизуйся по кнопке ниже!",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Залогиниться",
						Style:    discordgo.LinkButton,
						Disabled: false,
						URL:      "https://discordapp.com/api/oauth2/authorize?client_id=724098984389098368&permissions=0&scope=bot",
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
