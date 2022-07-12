package discord

import "github.com/bwmarrin/discordgo"

func (d *Discord) commands() map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate) {
	return map[string]func(s *discordgo.Session, i *discordgo.InteractionCreate){
		"help": d.commandHelp,
	}
}

func (d *Discord) commandHelp(s *discordgo.Session, i *discordgo.InteractionCreate) {
	_, err := d.ds.ChannelMessageSend(i.ChannelID, "Бот для администрирования сервера Rimas, функционал доступен только администраторам")
	if err != nil {
		d.logger.Errorf("commandHelp(): cant send message %s", err.Error())
		return
	}
}
