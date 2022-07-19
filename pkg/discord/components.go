package discord

import "github.com/bwmarrin/discordgo"

func (d *Discord) components() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"login":       d.componentLogin,
		"how_to_play": d.howToPlay,
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
	response := &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "ШАГ 1: Купить Arma 3",
							Style: discordgo.LinkButton,
							URL:   "https://store.steampowered.com/app/107410/Arma_3/?l=russian",
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "ШАГ 2: Подпишитесь на мод",
							Style: discordgo.LinkButton,
							URL:   "https://steamcommunity.com/sharedfiles/filedetails/?id=1368860933",
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "ШАГ 3: Установите TeamSpeak 3",
							Style: discordgo.LinkButton,
							URL:   "https://files.teamspeak-services.com/releases/client/3.5.6/TeamSpeak3-Client-win64-3.5.6.exe",
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "ШАГ 4: Установите плагин для TeamSpeak 3",
							Style: discordgo.LinkButton,
							URL:   "https://rimasrp.life/task_force_radio.ts3_plugin",
						},
					},
				},
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Видео инструкция",
							Style: discordgo.LinkButton,
							URL:   "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
						},
					},
				},
			},
			Embeds: []*discordgo.MessageEmbed{
				{
					Type: discordgo.EmbedTypeVideo,
					Video: &discordgo.MessageEmbedVideo{
						URL: "https://www.youtube.com/watch?v=dQw4w9WgXcQ",
					},
					Description: "Как начать играть?\nСледуй инструкции ниже\n" +
						"**Шаг 1**\nКупи и скачай ArmA 3 в Steam.\nhttps://store.steampowered.com/app/107410/Arma_3/?l=russian\n" +
						"**Шаг 2**\nПодпишись на мод Rimas Role Play в мастерской Steam.\nhttps://steamcommunity.com/sharedfiles/filedetails/?id=1368860933\n" +
						"**Шаг 3**\nСкачай клиент TeamSpeak и установи его.\nhttps://files.teamspeak-services.com/releases/client/3.5.6/TeamSpeak3-Client-win64-3.5.6.exe\n" +
						"**Шаг 4**\nСкачай плагин Task Force Radio и установи его.\nhttps://rimasrp.life/task_force_radio.ts3_plugin\nЗапуск\nЗапустите ArmA 3 в Steam, кликнув на кнопку играть.\n\n" +
						"В пункте **\"Моды\"** проверь, включен ли мод **Rimas Role Play**, если отключен — включи его.\n\n" +
						"Нажми на оранжевую кнопку играть в лаунчере ArmA 3.\n\n" +
						"Зайди в браузер серверов, нажми прямое подключение и введи \n**IP:** S1.RIMASRP.LIFE\n**Порт:** 2302\n\n" +
						"Удачи!",
					Title: "Как начать играть?",
					Color: 0x42FF00,
				},
			},
		},
	}
	err := s.InteractionRespond(i.Interaction, response)
	if err != nil {
		panic(err)
	}
}
