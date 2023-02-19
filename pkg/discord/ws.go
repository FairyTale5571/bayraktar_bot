package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
)

func (d *Discord) isExistUpdate(id string) bool {

	rows, err := d.db.Query("SELECT id FROM steam_updates WHERE id = ?", id)
	if err != nil {
		d.logger.Errorf("isExistUpdate(): Error getting update: %s", err.Error())
		return false
	}
	defer rows.Close()
	if rows.Next() {
		return true
	}
	return false
}

func (d *Discord) addUpdate(id, update string) {
	_, err := d.db.Exec("INSERT INTO `steam_updates` (`id`,`update`,`datetime`) VALUES (?,?,now())", id, update)
	if err != nil {
		d.logger.Errorf("addUpdate(): Error adding update: %s", err.Error())
		return
	}
}

func (d *Discord) checkUpdate() {
	text, id := d.steam.GetLatestUpdate("1368860933")
	if d.isExistUpdate(id) {
		return
	}
	d.addUpdate(id, text)
	if text == "" {
		d.logger.Warnf("checkUpdate(): Empty update")
		return
	}
	d.printUpdate(text)
}

const channelUpdates = "864647029049655337"

func (d *Discord) printUpdate(update string) {
	data := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeImage,
		Title:       "⚠️ Мод обновлен!",
		URL:         links.UrlMod,
		Description: "Обновите мод что-бы избежать проблем со входом на сервер!",
		Color:       0xFBFF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Внесенные изменения:", Value: update},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: d.cfg.URL + "/assets/images/logo.png",
		},
		Timestamp: "",
	}
	_, err := d.ds.ChannelMessageSendEmbed(channelUpdates, data)
	if err != nil {
		d.logger.Errorf("printUpdate(): Error sending message: %s", err.Error())
		return
	}
}
