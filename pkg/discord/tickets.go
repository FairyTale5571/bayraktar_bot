package discord

import (
	"fmt"
	"github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"

	"github.com/bwmarrin/discordgo"
)

const (
	parentActual  = "955941824516747294"
	parentArchive = "1012393930450534450"
)

func (d *Discord) printCreateTicket(channelID string) {
	embed := &discordgo.MessageEmbed{
		Title:       "Помощь",
		Description: "Если вы столкнулись с проблемой, нажмите кнопку \"Открыть тикет\"",
		Color:       0x00ff00,
	}
	msg := &discordgo.MessageSend{
		Embed: embed,
		TTS:   false,
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

func (d *Discord) saveTicket(creatorId string) (int64, error) {
	rows, err := d.db.Exec("INSERT INTO discord_tickets (creator_id, insert_date) VALUES (?, now())", creatorId)
	if err != nil {
		return 0, err
	}
	lId, err := rows.LastInsertId()
	if err != nil {
		return 0, err
	}
	return lId, nil
}

func (d *Discord) updateTicket(id int64, channelId string) {
	_, err := d.db.Exec("UPDATE discord_tickets SET channel_id = ? WHERE id = ?", channelId, id)
	if err != nil {
		d.logger.Errorf("cant update ticket: %v", err)
		return
	}
}

func (d *Discord) getTicketOwner(channelID string) (string, error) {
	rows, err := d.db.Query("SELECT creator_id FROM discord_tickets WHERE channel_id = ?", channelID)
	defer rows.Close()
	if err != nil {
		return "", err
	}
	var creatorId string
	for rows.Next() {
		err := rows.Scan(&creatorId)
		if err != nil {
			return "", err
		}
	}
	return creatorId, nil
}

func (d *Discord) isTicketOpened(userID string) bool {
	rows, err := d.db.Query("SELECT id FROM discord_tickets WHERE creator_id = ?", userID)
	defer rows.Close()
	if err != nil {
		return false
	}
	if rows.Next() {
		return true
	}
	return false
}

func (d *Discord) createTickets(guildID string, user *discordgo.User) (*discordgo.Channel, error) {
	if d.isTicketOpened(user.ID) {
		return nil, errorUtils.ErrTicketOpened
	}
	ticketId, err := d.saveTicket(user.ID)
	if err != nil {
		return nil, err
	}
	data := discordgo.GuildChannelCreateData{
		Name:                 fmt.Sprintf("support-%d", ticketId),
		Type:                 discordgo.ChannelTypeGuildText,
		PermissionOverwrites: nil,
		ParentID:             parentActual,
	}
	channel, err := d.ds.GuildChannelCreateComplex(guildID, data)
	if err != nil {
		d.logger.Errorf("cant create channel: %v", err)
		return nil, err
	}
	d.updateTicket(ticketId, channel.ID)
	_, err = d.ds.ChannelMessageSendComplex(channel.ID, &discordgo.MessageSend{
		Content: fmt.Sprintf("Здравствуйте, %s!\n", user.Mention()) +
			"Опишите подробно проблему с которой вы столкнулись, приложите скриншоты, последовательность ваших действий которые привели к проблеме\n" +
			"Мы постараемся решить вашу проблему как можно скорее\n" +
			"Если проблема не актуальна, нажмите кнопку \"Закрыть тикет\"",
		Components: []discordgo.MessageComponent{
			discordgo.ActionsRow{
				Components: []discordgo.MessageComponent{
					discordgo.Button{
						Label:    "Закрыть тикет",
						CustomID: "close_ticket",
						Style:    discordgo.DangerButton,
					},
				},
			},
		},
	})
	if err != nil {
		d.logger.Errorf("cant send message: %v", err)
		return nil, err
	}
	err = d.ds.ChannelPermissionSet(channel.ID, user.ID, discordgo.PermissionOverwriteTypeMember, 117824, 0)
	if err != nil {
		d.logger.Errorf("cant set permission: %v", err)
		return nil, err
	}
	return channel, nil
}

func (d *Discord) closeTicket(s *discordgo.Session, i *discordgo.InteractionCreate) {
	err := d.ds.InteractionRespond(i.Interaction, &discordgo.InteractionResponse{
		Type: discordgo.InteractionResponseChannelMessageWithSource,
		Data: &discordgo.InteractionResponseData{
			Content: fmt.Sprintf("%s закрыл тикет", i.Interaction.Member.User.Mention()),
		},
	})
	if err != nil {
		d.logger.Errorf("cant respond interaction: %v", err)
	}
	_ = d.ds.ChannelMessageDelete(i.Interaction.ChannelID, i.Interaction.Message.ID)

	channelID := i.Interaction.ChannelID
	_, err = d.db.Exec("DELETE FROM discord_tickets WHERE channel_id = ?", channelID)
	if err != nil {
		d.logger.Errorf("cant delete ticket from database: %v", err)
		return
	}
	_, err = d.ds.ChannelDelete(channelID)
	if err != nil {
		d.logger.Errorf("cant delete channel: %v", err)
		return
	}
}

func (d *Discord) createChannel(guildID string) {
	channel, err := d.ds.GuildChannelCreate(guildID, "support", discordgo.ChannelTypeGuildText)
	if err != nil {
		d.logger.Errorf("cant create channel: %v", err)
		return
	}
	d.logger.Infof("channel created: %s", channel.ID)
}
