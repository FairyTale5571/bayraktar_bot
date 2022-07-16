package discord

import (
	"fmt"
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help":       d.commandHelp,
		"copy-role":  d.commandCopyRole,
		"give-boost": d.commandGiveBoost,
		"get-him":    d.commandGetHim,
	}
}

func (d *Discord) commandHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := d.ds.ChannelMessageSend(i.ChannelID, "Бот для администрирования сервера Rimas, функционал доступен только администраторам")
	if err != nil {
		d.logger.Errorf("commandHelp(): cant send message %s", err.Error())
		return
	}
}

func (d *Discord) commandCopyRole(s *discordgo.Session, i *discordgo.InteractionCreate) {
	roleId := i.ApplicationCommandData().Options[0].RoleValue(s, "").ID
	role, err := d.findRole(i.Interaction.GuildID, roleId)
	if err != nil {
		d.printHiddenMessageInteraction(i, "ошибка при создании роли "+err.Error())
		return
	}
	nameNewRole := i.ApplicationCommandData().Options[1].StringValue()
	newRole, err := s.GuildRoleCreate(i.Interaction.GuildID)
	if err != nil {
		d.printHiddenMessageInteraction(i, "ошибка при создании роли "+err.Error())
		return
	}
	newRole, err = s.GuildRoleEdit(i.Interaction.GuildID, newRole.ID, nameNewRole, role.Color, role.Hoist, role.Permissions, role.Mentionable)
	if err != nil {
		d.printHiddenMessageInteraction(i, "ошибка при изменении роли "+err.Error())
		return
	}
	d.printHiddenMessageInteraction(i, "Роль "+role.Mention()+" скопирована в "+newRole.Mention())
}

func (d *Discord) commandGiveBoost(s *discordgo.Session, i *discordgo.InteractionCreate) {
	user := i.ApplicationCommandData().Options[0].UserValue(nil)
	sender := i.Interaction.Member.User
	if !d.isAdmin(i.Interaction.ChannelID, sender.ID) {
		d.printHiddenMessageInteraction(i, "Вы не администратор")
		return
	}
	d.giveBoostPresent(i.ChannelID, user)
	d.printHiddenMessageInteraction(i, "Запрос отправлен")
}

func (d *Discord) commandGetHim(s *discordgo.Session, i *discordgo.InteractionCreate) {
	var err error
	user := i.ApplicationCommandData().Options[0].UserValue(d.ds)
	sender := i.Interaction.Member.User
	if !d.isAdmin(i.Interaction.ChannelID, sender.ID) {
		d.logger.Errorf("commandGetHim(): %s is not admin", sender.ID)
		d.printHiddenMessageInteraction(i, "Вы не администратор")
		return
	}
	_id, err := d.getUserSteamId(user.ID)
	if err != nil {
		d.logger.Errorf("commandGetHim(): %s", err.Error())
		d.printHiddenMessageInteraction(i, "Ошибка при получении SteamID")
		return
	}

	_player := d.getPlayerInformation(_id)
	if _player == nil {
		d.logger.Errorf("commandGetHim(): %s", err.Error())
		d.printHiddenMessageInteraction(i, "Ошибка при получении информации о игроке")
		return
	}

	d.printHiddenEmbedInteraction(i, &discordgo.MessageEmbed{
		Title: "Пользователь",
		Fields: []*discordgo.MessageEmbedField{
			{Name: "ID:", Value: fmt.Sprintf("%d", _player.Id), Inline: true},
			{Name: "Steam:", Value: _player.Uid, Inline: true},
			{Name: "Профиль:", Value: _player.Name, Inline: true},
			{Name: "ФНИ:", Value: fmt.Sprintf("%s '%s' %s", _player.FirstName.String, _player.NickName.String, _player.LastName.String)},
			{Name: "Дата регистрации:", Value: _player.InsertTime.Format("02.01.2006 15:04:05")},
			{Name: "Дата последнего входа:", Value: _player.LastConnected.Format("02.01.2006 15:04:05")},
			{Name: "Всего времени в игре:", Value: fmt.Sprintf("%s", secondsToDate(_player.TotalTime))},

			{Name: "Наличных:", Value: fmt.Sprintf("$%d", _player.Cash), Inline: true},
			{Name: "В банке:", Value: fmt.Sprintf("$%d", _player.Bank), Inline: true},
			{Name: "RC:", Value: fmt.Sprintf("%d", _player.RC), Inline: true},
		},
		Color: 0x00ff00,
	})
}
