package discord

import "github.com/bwmarrin/discordgo"

func (d *Discord) printCreateTicket(channelID string) {

	embed := &discordgo.MessageEmbed{
		Title:       "Помощь",
		Description: "Если вы столкнулись с проблемой, нажмите кнопку \"Открыть тикет\"",
		Color:       0x00ff00,
	}
	msg := &discordgo.MessageSend{
		Embed: embed,
		TTS:   true,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Открыть тикет",
						CustomID: "create_ticket",
						Style:    discordgo.SuccessButton,
					},
				},
			},
		},
	}
	_, err := d.ds.ChannelMessageSendComplex(channelID, msg)
	if err != nil {
		d.logger.Errorf("cant send message: %v", err)
	}

}
