package discord

import (
	"fmt"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"time"

	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"
)

const (
	parentActual   = "955941824516747294"
	channelReports = "1016722224461402153"
)

func (d *Discord) getTicketCreateDate(channelID string) time.Time {
	var res time.Time
	rows := d.db.QueryRow("SELECT insert_date FROM discord_tickets WHERE channel_id = ?", channelID)
	if err := rows.Scan(&res); err != nil {
		d.logger.Errorf("cant scan ticket create date: %v", err)
		return res
	}
	return res
}

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

func (d *Discord) checkAuthorize(member *discordgo.Member) {
	for _, role := range member.Roles {
		if role == d.cfg.RegRoleID {
			return
		}
	}
	ch, err := d.ds.UserChannelCreate(member.User.ID)
	if err != nil {
		d.logger.Errorf("cant create channel: %v", err)
		return
	}
	_, err = d.ds.ChannelMessageSendComplex(ch.ID, &discordgo.MessageSend{
		Embed: &discordgo.MessageEmbed{
			Color: 0xFFC500,
			Title: "Необходима регистрация",
			Description: "Здравствуйте! Вы получили это сообщение, так как вы создали тикет, но все еще не авторизованы на сервере.\n" +
				"Для того чтобы мы могли вам помочь как можно быстрее, пожалуйста, пройдите авторизацию в канале <#872192873825701968>.\n" +
				"**Если вы не заходили на сервер и у вас проблема с входом, то проигнорируйте данное сообщение**",
		},
	})
	if err != nil {
		d.logger.Errorf("cant send message: %v", err)
		return
	}
}

func (d *Discord) createTickets(guildID string, user *discordgo.Member) (*discordgo.Channel, error) {
	if d.isTicketOpened(user.User.ID) {
		return nil, errorUtils.ErrTicketOpened
	}
	d.checkAuthorize(user)
	ticketId, err := d.saveTicket(user.User.ID)
	if err != nil {
		return nil, err
	}
	data := discordgo.GuildChannelCreateData{
		Name:                 fmt.Sprintf("%s-%d", user.User.Username, ticketId),
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
	err = d.ds.ChannelPermissionSet(channel.ID, user.User.ID, discordgo.PermissionOverwriteTypeMember, 117824, 0)
	if err != nil {
		d.logger.Errorf("cant set permission: %v", err)
		return nil, err
	}
	return channel, nil
}

func (d *Discord) serializeReport(channelID, closerID string) (*models.TicketReport, error) {
	var res models.TicketReport
	channel, err := d.ds.Channel(channelID)
	if err != nil {
		d.logger.Errorf("cant get channel: %v", err)
		return nil, err
	}
	res.ChannelName = channel.Name
	res.ChannelID = channel.ID
	ticketOwner, err := d.getTicketOwner(channel.ID)
	if err != nil {
		d.logger.Errorf("cant get ticket owner: %v", err)
		return nil, err
	}
	res.AuthorID = ticketOwner
	res.ClosedBy = closerID
	res.ClosedAt = time.Now()
	res.OpenedAt = d.getTicketCreateDate(channel.ID)

	messages, err := d.ds.ChannelMessages(channel.ID, 100, "", "", "")
	if err != nil {
		d.logger.Errorf("cant get messages: %v", err)
		return nil, err
	}
	res.Messages = messages

	return &res, nil
}

func (d *Discord) sendReport(report *models.TicketReport) {
	embed := discordgo.MessageEmbed{
		Title: "Отчет по тикету",
		Description: fmt.Sprintf("Автор тикета: <@%s>\nЗакрыл тикет: <@%s>\nВремя открытия: %s\nВремя закрытия: %s\n**Содержание**:\n",
			report.AuthorID,
			report.ClosedBy,
			report.OpenedAt.Format("15:04:05 02-01-2006"),
			report.ClosedAt.Format("15:04:05 02-01-2006"),
		),
		Color: 0x00ff00,
	}
	for i := len(report.Messages) - 1; i >= 0; i-- {
		if report.Messages[i].Content == "" {
			continue
		}
		embed.Fields = append(embed.Fields, &discordgo.MessageEmbedField{
			Name:   fmt.Sprintf("%s: %s", report.Messages[i].Author.Username, report.Messages[i].Timestamp.Format("02.01.2006 15:04:05")),
			Value:  report.Messages[i].Content,
			Inline: false,
		})
	}
	// send to channel with reports
	_, err := d.ds.ChannelMessageSendEmbed(channelReports, &embed)
	if err != nil {
		d.logger.Errorf("cant send message: %v", err)
	}
	// send to user
	userChannel, err := d.ds.UserChannelCreate(report.AuthorID)
	if err != nil {
		d.logger.Errorf("cant create channel: %v", err)
		return
	}
	_, err = d.ds.ChannelMessageSendComplex(userChannel.ID, &discordgo.MessageSend{
		Embed: &embed,
		/*
			Components: []discordgo.MessageComponent{
				discordgo.ActionsRow{
					Components: []discordgo.MessageComponent{
						discordgo.Button{
							Label: "Мне понравилось обслуживание",
							Emoji: discordgo.ComponentEmoji{
								Name: "👍",
							},
							CustomID: "problem_solved",
							Style:    discordgo.SuccessButton,
						},
						discordgo.Button{
							Label: "Обслуживание не понравилось",
							Emoji: discordgo.ComponentEmoji{
								Name: "👎",
							},
							CustomID: "problem_not_solved",
							Style:    discordgo.DangerButton,
						},
					},
				},
			},

		*/
	})
	if err != nil {
		d.logger.Errorf("cant send message: %v", err)
	}
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
	if report, err := d.serializeReport(i.ChannelID, i.Member.User.ID); err == nil {
		d.sendReport(report)
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
