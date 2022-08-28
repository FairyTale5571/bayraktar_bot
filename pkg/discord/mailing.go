package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
)

func (d *Discord) madeFields(fields []struct {
	Name  string `json:"name"`
	Value string `json:"value"`
},
) []*discordgo.MessageEmbedField {
	var res []*discordgo.MessageEmbedField
	for _, v := range fields {
		res = append(res, &discordgo.MessageEmbedField{
			Name:  v.Name,
			Value: v.Value,
		})
	}
	return res
}

func (d *Discord) serializeEmbed(embeds models.Embed) *discordgo.MessageEmbed {
	return &discordgo.MessageEmbed{
		Title:       embeds.Title,
		Description: embeds.Description,
		Color:       embeds.Color,
		Footer: &discordgo.MessageEmbedFooter{
			Text:    embeds.Footer.Text,
			IconURL: embeds.Footer.IconUrl,
		},
		Fields:    d.madeFields(embeds.Fields),
		URL:       embeds.Url,
		Timestamp: embeds.Timestamp,
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: embeds.Thumbnail.Url,
		},
	}
}

func (d *Discord) serializeEmbeds(embeds models.Embeds) []*discordgo.MessageEmbed {
	var res []*discordgo.MessageEmbed
	for _, v := range embeds.Embeds {
		res = append(res, d.serializeEmbed(v))
	}
	return res
}

func (d *Discord) makeButtonLink() []discordgo.MessageComponent {
	return []discordgo.MessageComponent{
		discordgo.ActionsRow{
			Components: []discordgo.MessageComponent{
				discordgo.Button{
					URL:   links.UrlDiscord,
					Label: "Отправлено с сервера Rocket Life",
					Style: discordgo.LinkButton,
				},
			},
		},
	}
}

func (d *Discord) SendMassive(guildID string, embeds models.Embeds) {
	data := &discordgo.MessageSend{
		Embeds:     d.serializeEmbeds(embeds),
		Components: d.makeButtonLink(),
	}
	members, err := d.getAllMembers(guildID)
	if err != nil {
		d.logger.Errorf("Error while getting members: %s", err)
		return
	}
	for _, member := range members {
		if member.User.Bot {
			continue
		}
		d.sendPrivateMessage(member.User.ID, data)
	}
}

func (d *Discord) SendDirect(userID string, embeds models.Embeds) {
	for _, v := range embeds.Embeds {
		embed := d.serializeEmbed(v)
		data := &discordgo.MessageSend{
			Embed:      embed,
			Components: d.makeButtonLink(),
		}
		d.sendPrivateMessage(userID, data)
	}
}

func (d *Discord) SendToChannel(guildID, channelID string, embeds models.Embeds) {
	for _, v := range embeds.Embeds {
		embed := d.serializeEmbed(v)
		data := &discordgo.MessageSend{
			Content: embeds.Content,
			Embed:   embed,
		}
		_, err := d.ds.ChannelMessageSendComplex(channelID, data)
		if err != nil {
			d.logger.Errorf("SendToChannel(): Error while sending message: %s", err)
		}
	}
}
