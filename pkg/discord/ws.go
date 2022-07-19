package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage"
)

func (d *Discord) checkUpdate() {
	var err error
	text, id := d.steam.GetLatestUpdate("1368860933")

	lastUpdate, _ := d.rdb.Get("lastUpdate_1368860933", storage.LastWsUpdate)
	if lastUpdate == id {
		return
	}
	err = d.rdb.Set("lastUpdate_1368860933", id, storage.LastWsUpdate)
	if err != nil {
		d.logger.Errorf("checkUpdate(): Error setting last update: %s", err.Error())
		return
	}
	if text == "" {
		text = "Незначительные изменения"
	}
	d.printUpdate(text)
}

func (d *Discord) printUpdate(update string) {
	data := &discordgo.MessageEmbed{
		Type:        discordgo.EmbedTypeImage,
		Title:       "⚠️ Мод обновлен!",
		URL:         "https://steamcommunity.com/sharedfiles/filedetails/?id=1368860933",
		Description: "Обновите мод что-бы избежать проблем со входом на сервер!",
		Color:       0xFBFF00,
		Fields: []*discordgo.MessageEmbedField{
			{Name: "Внесенные изменения:", Value: update},
		},
		Thumbnail: &discordgo.MessageEmbedThumbnail{
			URL: "https://i.imgur.com/fkuPE2b.gif",
		},
		Timestamp: "",
	}
	_, err := d.ds.ChannelMessageSendEmbed("864647029049655337", data)
	if err != nil {
		d.logger.Errorf("printUpdate(): Error sending message: %s", err.Error())
		return
	}
}
