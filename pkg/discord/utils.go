package discord

import (
	"fmt"
	"math/rand"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"
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
	defer rows.Close() // nolint: errcheck

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
	defer rows.Close() // nolint: errcheck

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
			g.name, p.group_level,
			p.insert_time, p.last_connected, p.total_time
		FROM players p
		INNER JOIN groups g ON g.id = p.group_id
		WHERE p.playerid = ?`, steamId)
	defer rows.Close() // nolint: errcheck
	if err != nil {
		d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
		return nil
	}
	for rows.Next() {
		if err := rows.Scan(&_player.Id, &_player.Uid, &_player.Name, &_player.NickName, &_player.FirstName, &_player.LastName, &_player.Cash, &_player.Bank, &_player.RC, &_player.GroupName, &_player.GroupLevel, &_player.InsertTime, &_player.LastConnected, &_player.TotalTime); err != nil {
			d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
			return nil
		}
	}
	if _player.GroupName.Valid {
		rows, err = d.db.Query(`SELECT JSON_EXTRACT(titles, '$[?][1]') from groups where name = ?`, _player.GroupName.String)
		defer rows.Close() // nolint: errcheck
		if err != nil {
			d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
			return nil
		}
		for rows.Next() {
			if err := rows.Scan(&_player.GroupLevelName); err != nil {
				d.logger.Errorf("getPlayerInformation(): cant get player data %s", err.Error())
				return nil
			}
		}
	}

	return &_player
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
