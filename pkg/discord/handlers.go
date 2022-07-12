package discord

import (
	"github.com/bwmarrin/discordgo"
)

func (d *Discord) onBotUp(s *discordgo.Session, r *discordgo.Ready) {
	d.logger.Infof("Bot %s is Up", r.User.Username)
	return
}

func (d *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	return
}

func (d *Discord) onUserChanged(s *discordgo.Session, i *discordgo.GuildMemberUpdate) {
	return
}

func (d *Discord) onUserConnected(s *discordgo.Session, i *discordgo.GuildMemberAdd) {
	return
}

func (d *Discord) onUserDisconnected(s *discordgo.Session, i *discordgo.MessageReactionAdd) {
	return
}

func (d *Discord) onGuildCreate(s *discordgo.Session, g *discordgo.GuildCreate) {
	for _, v := range commands {
		_, err := d.ds.ApplicationCommandCreate(d.ds.State.User.ID, g.ID, v)
		if err != nil {
			d.logger.Errorf("cant create command %s", err.Error())
			continue
		}
		d.logger.Infof("command %s created on server %s", v.Name, g.Name)
	}
}

func (d *Discord) onCommandsCall(s *discordgo.Session, i *discordgo.InteractionCreate) {

	switch i.Type {
	case discordgo.InteractionApplicationCommand:
		if h, ok := d.commands()[i.ApplicationCommandData().Name]; ok {
			h(s, i)
		}
	}
}
