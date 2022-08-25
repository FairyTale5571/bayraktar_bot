package discord

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
)

func (d *Discord) getAllMembers(guildId string) ([]*discordgo.Member, error) {
	var members []*discordgo.Member
	after := ""
	for {
		users, err := d.ds.GuildMembers(guildId, after, 1000)
		if err != nil {
			d.logger.Errorf("get users error: %s\n", err.Error())
			break
		}
		members = append(members, users...)
		after = users[len(users)-1].User.ID
		if len(users) != 1000 {
			break
		}
	}
	return members, nil
}

func (d *Discord) isAdmin(channelId, userId string) bool {
	permission, err := d.ds.UserChannelPermissions(userId, channelId)
	if err != nil {
		d.logger.Errorf("isAdmin(): cant get user permissions %s", err.Error())
		return false
	}
	return permission&discordgo.PermissionAdministrator != 0
}

func (d *Discord) findRole(guildId, roleId string) (*discordgo.Role, error) {
	roles, err := d.ds.GuildRoles(guildId)
	if err != nil {
		return &discordgo.Role{}, err
	}
	for _, elem := range roles {
		if elem.ID == roleId {
			return elem, nil
		}
	}
	return &discordgo.Role{}, errorUtils.ErrRoleNotFound
}

func (d *Discord) printHiddenMessageInteraction(i *discordgo.InteractionCreate, msg string) {
	err := d.ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:   1 << 6,
			Content: msg,
		},
	})
	if err != nil {
		d.logger.Errorf("printHiddenMessageInteraction(): cant send message %s", err.Error())
		return
	}
}

func (d *Discord) printHiddenEmbedInteraction(i *discordgo.InteractionCreate, msg *discordgo.MessageEmbed) {
	err := d.ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Flags:  1 << 6,
			Embeds: []*discordgo.MessageEmbed{msg},
		},
	})
	if err != nil {
		d.logger.Errorf("printHiddenMessageInteraction(): cant send message %s", err.Error())
		return
	}
}

func (d *Discord) getUserSteamId(userId string) (string, error) {
	rows, err := d.db.Query("SELECT uid FROM discord_users WHERE discord_uid = ?", userId)
	if err != nil {
		return "", err
	}
	defer rows.Close() // nolint: not needed

	var uid string
	for rows.Next() {
		if err := rows.Scan(&uid); err != nil {
			return "", err
		}
	}
	if uid == "" {
		return "", errorUtils.ErrSteamUserNotFound
	}
	return uid, nil
}

func (d *Discord) getRandomVehicle() *Vehicles {
	var veh Vehicles
	rows, err := d.db.Query("SELECT d.classname, d.image, m.displayName FROM discord_boosters d " +
		"INNER JOIN lk_mapobjects m ON m.classname = d.classname " +
		"WHERE active = 1 " +
		"ORDER BY RAND() " +
		"LIMIT 1")
	defer rows.Close() // nolint: not needed

	if err != nil {
		d.logger.Errorf("getRandomVehicle(): cant get random vehicle %s", err.Error())
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&veh.Classname, &veh.Image, &veh.DisplayName); err != nil {
			d.logger.Errorf("getRandomVehicle(): cant get random vehicle %s", err.Error())
			return nil
		}
	}
	return &veh
}

func (d *Discord) giveBoostPresent(channelId string, user *discordgo.User) {
	var err error
	player, err := d.getUserSteamId(user.ID)
	if err != nil {
		d.logger.Errorf("giveBoostPresent(): cant get user steam id %s", err.Error())
		_, err = d.ds.ChannelMessageSend(channelId, user.Mention()+"\nМы не нашли ваш аккаунт на сервере, привяжите ваш аккаунт и напишите администрации за получением бонуса")
		if err != nil {
			d.logger.Errorf("giveBoostPresent(): cant send message %s", err.Error())
			return
		}
		return
	}
	vehicle := d.getRandomVehicle()
	_, err = d.db.Exec("INSERT INTO vehicles SET servermap = 'RRpMap',classname = ?, pid = ?, plate = ?,"+
		"type = 'Car', alive = '1', active = '0', inventory = '[[],0]',color = 'default', material = 'default', gear = '[]', damage = '0', hitpoints = '[]', baseprice = 10000, spname = 'none', parking = '[]', maxslots = 60, tuning_data = '[[\"nitro\"],[\"tracker\"],[\"breaking\"],[\"seatbelt\"]]', distance = '0', deleted_at = NULL, comment = ''",
		player, vehicle.Classname, generatePlateNumber())

	if err != nil {
		d.logger.Errorf("giveBoostPresent(): cant insert vehicle %s", err.Error())
		_, err = d.ds.ChannelMessageSend(channelId, user.Mention()+"\nНе удалось добавить автомобиль в базу данных, обратитесь к администрации")
		if err != nil {
			d.logger.Errorf("giveBoostPresent(): cant send message %s", err.Error())
			return
		}
		return
	}

	_, _ = d.ds.ChannelMessageSend(channelId, user.Mention())
	_, _ = d.ds.ChannelMessageSendEmbed(channelId, &discordgo.MessageEmbed{
		URL:         "",
		Type:        discordgo.EmbedTypeImage,
		Title:       "Nitro Booster",
		Description: fmt.Sprintf("Спасибо за буст сервера!\nТвой подарок %v уже доступен на сервере!", vehicle.DisplayName),
		Timestamp:   "",
		Color:       0x9300FF,
		Footer: &discordgo.MessageEmbedFooter{
			Text:         "Nitro Boost",
			IconURL:      "",
			ProxyIconURL: "",
		},
		Image: &discordgo.MessageEmbedImage{
			URL:      vehicle.Image,
			ProxyURL: "",
			Width:    0,
			Height:   0,
		},
	})
}

func (d *Discord) getPlayerInformation(steamId string) *PlayerData {
	var _player PlayerData
	//)
	rows, err := d.db.Query(`SELECT p.uid, p.playerid,
			p.name, p.nick_name, p.first_name, p.last_name,
			p.cash, p.bankacc, p.EPoint,
			p.group_id, p.group_level,
			p.insert_time, p.last_connected, p.total_time
		FROM players p
		WHERE p.playerid = ?`, steamId)
	defer rows.Close() // nolint: not needed
	if err != nil {
		d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&_player.Id, &_player.Uid, &_player.Name, &_player.NickName, &_player.FirstName, &_player.LastName, &_player.Cash, &_player.Bank, &_player.RC, &_player.GroupID, &_player.GroupLevel, &_player.InsertTime, &_player.LastConnected, &_player.TotalTime); err != nil {
			d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
			return nil
		}
	}
	if _player.GroupID != -1 {
		rows, err = d.db.Query(`SELECT name, JSON_EXTRACT(titles, '$[?][1]') from groups where id = ?`, _player.GroupID)
		defer rows.Close() // nolint: not needed
		if err != nil {
			d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
			return nil
		}
		for rows.Next() {
			if err := rows.Scan(&_player.GroupName, &_player.GroupLevelName); err != nil {
				d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
				return nil
			}
		}
	}

	return &_player
}

func (d *Discord) printPrivateMessage(userId, text string) {
	channel, err := d.ds.UserChannelCreate(userId)
	if err != nil {
		d.logger.Errorf("printPrivateMessage(): cant create channel %s", err.Error())
		return
	}
	_, err = d.ds.ChannelMessageSend(channel.ID, text)
	if err != nil {
		d.logger.Errorf("printPrivateMessage(): cant send message %s", err.Error())
		return
	}
}

func (d *Discord) RegisterUser(guildId, userId, steamId string) {
	var _userId string
	_ = d.db.QueryRow("SELECT id FROM discord_users WHERE uid = ? limit 1", steamId).Scan(&_userId)
	if _userId != "" {
		d.printPrivateMessage(userId, "Этот Steam аккаунт уже зарегистрирован")
		return
	}
	_ = d.db.QueryRow("SELECT uid FROM players WHERE playerid = ? limit 1", steamId).Scan(&_userId)
	if _userId == "" {
		d.printPrivateMessage(userId, "Мы не нашли ваш Steam аккаунт в базе данных\nВероятно это из-за того, что вы еще не играли на сервере")
		return
	}
	user, err := d.ds.GuildMember(guildId, userId)
	if err != nil {
		d.logger.Errorf("RegisterUser(): cant get user %s", err.Error())
		return
	}
	_, err = d.db.Exec("INSERT INTO discord_users (uid, discord_uid, discord_name, discord_discriminator) VALUES (?, ?, ?, ?)", steamId, user.User.ID, user.User.Username, user.User.Discriminator)
	if err != nil {
		d.logger.Errorf("RegisterUser(): cant insert user %s", err.Error())
		return
	}
	err = d.ds.GuildMemberRoleAdd(guildId, userId, "864630308242849825")
	if err != nil {
		d.logger.Errorf("RegisterUser(): cant add role %s", err.Error())
		return
	}
	d.printPrivateMessage(userId, "Вы успешно зарегистрированы в сервере!\nДоступ к каналам предоставлен!")
}

func (d *Discord) deleteUser(userId string) {
	_, err := d.db.Exec("DELETE FROM discord_users WHERE discord_uid = ?", userId)
	if err != nil {
		d.logger.Errorf("deleteUser(): cant delete user %s", err.Error())
	}
}

func (d *Discord) printWelcome(userID, guildID string) {
	guild, err := d.ds.Guild(guildID)
	if err != nil {
		d.logger.Errorf("printWelcome(): cant get guild %s", err.Error())
		return
	}
	channel, err := d.ds.UserChannelCreate(userID)
	if err != nil {
		d.logger.Errorf("printWelcome(): cant create channel %s", err.Error())
		return
	}
	embed := &discordgo.MessageEmbed{
		Type:  discordgo.EmbedTypeImage,
		Title: "Добро пожаловать на сервер **" + guild.Name + "**!",
		Description: "Чтобы получить полный доступ к серверу, тебе нужно привязать твой аккаунт к нашему серверу!\n" +
			"Сделать это можно по кнопке ниже!\n" +
			"Не забудь ознакомится с правилами поведения в игре\n\n" +
			"Прочувствуй атмосферу удивительного мира ролевой игры с реалистичным миром!",
		Color: 0x00ff00,
		Image: &discordgo.MessageEmbedImage{
			URL: d.steam.GetItemLogo("1368860933"),
		},
	}
	msg := &discordgo.MessageSend{
		Embed: embed,
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Привязать аккаунт",
						URL:   d.steam.GetAuthLink(guildID, userID),
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🔗",
						},
					},
					discordgo.Button{
						Label:    "Как начать играть?",
						Style:    discordgo.SuccessButton,
						CustomID: "how_to_play",
						Emoji: discordgo.ComponentEmoji{
							Name: "📚",
						},
					},
				},
			},
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label: "Форум",
						URL:   links.UrlForum,
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🔗",
						},
					},
					discordgo.Button{
						Label: "Личный кабинет",
						URL:   links.UrlLk,
						Style: discordgo.LinkButton,
						Emoji: discordgo.ComponentEmoji{
							Name: "🔗",
						},
					},
				},
			},
		},
	}
	_, err = d.ds.ChannelMessageSendComplex(channel.ID, msg)
	if err != nil {
		d.logger.Errorf("printWelcome(): cant send message %s", err.Error())
	}
}

var letters = []rune("ABEIKMHOPCTXZ")

func randStringRune(n int) string {
	b := make([]rune, n)
	rand.Seed(time.Now().UnixNano())
	for i := range b {
		b[i] = letters[rand.Intn(len(letters))]
	}
	return string(b)
}

func randInt() int {
	rand.Seed(time.Now().UnixNano())
	min := 1000
	max := 9999
	return rand.Intn(max-min+1) + min
}

func generatePlateNumber() string {
	return fmt.Sprintf("DS %d %v", randInt(), randStringRune(2))
}

func secondsToDate(seconds uint64) string {
	return fmt.Sprintf("%d дней %d часов %d минут", seconds/86400, (seconds%86400)/3600, (seconds%3600)/60)
}

func (d *Discord) getHow2Play() (*discordgo.MessageEmbed, []discordgo.MessageComponent) {
	components := []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "ШАГ 1: Купить Arma 3",
					Style: discordgo.LinkButton,
					URL:   links.UrlGame,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "ШАГ 2: Подпишитесь на мод",
					Style: discordgo.LinkButton,
					URL:   links.UrlMod,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "ШАГ 3: Установите TeamSpeak 3",
					Style: discordgo.LinkButton,
					URL:   links.UrlTeamspeak,
				},
			},
		},
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					Label: "ШАГ 4: Установите плагин для TeamSpeak 3",
					Style: discordgo.LinkButton,
					URL:   d.cfg.URL + "/assets/files/task_force_radio.ts3_plugin",
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
	}
	embed := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeImage,
		Title:       "Как начать играть",
		Description: "Следуйте инструкции ниже",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Шаг 1", Value: "Купите и скачайте ArmA 3 в Steam.\n" + links.UrlGame},
			{Name: "Шаг 2", Value: "Подпишитесь на мод Rocket Life в мастерской Steam.\n" + links.UrlMod},
			{Name: "Шаг 3", Value: "Скачайте клиент TeamSpeak и установите его.\n" + links.UrlTeamspeak},
			{Name: "Шаг 4", Value: "Скачайте плагин Task Force Radio и установите его.\n" + d.cfg.URL + "/assets/files/task_force_radio.ts3_plugin"},
			{Name: "Запуск", Value: "Запустите ArmA 3 в Steam, кликнув на кнопку играть.\n\nВ пункте \"Моды\" проверьте, включен ли мод **Rocket Life**, если отключен — включите его.\n\nНажмите на оранжевую кнопку играть в лаунчере ArmA 3.\n\nВ правом верхнем углы игры зайдите в свой профиль и укажите имя и фамилию вашего персонажа.\n\nЗайдите в браузер серверов и нажмите прямое подключение.\n\n" + links.UrlServer},
		},
		Color: 0x8700ff,
		Image: &discordgo.MessageEmbedImage{
			URL: d.cfg.URL + "/assets/images/big_logo.jpg",
		},
	}
	return embed, components
}

const (
	channelLog = "1011990018623025185"
)

func (d *Discord) printLog(logstr string) {
	_, err := d.ds.ChannelMessageSend(channelLog, logstr)
	if err != nil {
		d.logger.Errorf("printLog(): %s", err)
		return
	}
}

func (d *Discord) sendPrivateMessage(userID string, message *discordgo.MessageSend) {
	ch, err := d.ds.UserChannelCreate(userID)
	if err != nil {
		d.logger.Errorf("sendPrivateMessage(): Error user channel create: %s", err.Error())
		return
	}
	_, err = d.ds.ChannelMessageSendComplex(ch.ID, message)
	if err != nil {
		d.logger.Errorf("sendPrivateMessage(): Error sending message: %s", err.Error())
		return
	}
	d.logger.Infof("sendPrivateMessage(): Message sent to %s", userID)
}
