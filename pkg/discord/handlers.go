package discord

import (
	"strings"
	"time"

	"github.com/bwmarrin/discordgo"
)

func (d *Discord) onBotUp(s *discordgo.Session, r *discordgo.Ready) {
	d.logger.Infof("Bot %s is Up", r.User.Username)
	return
}

func (d *Discord) onMessageCreate(s *discordgo.Session, m *discordgo.MessageCreate) {
	// ignore messages from bots
	if m.Author.Bot {
		return
	}

	// catch discord boosters
	switch m.Message.Type {
	case
		discordgo.MessageTypeGuildMemberJoin:
	case
		discordgo.MessageTypeUserPremiumGuildSubscriptionTierOne,
		discordgo.MessageTypeUserPremiumGuildSubscriptionTierTwo,
		discordgo.MessageTypeUserPremiumGuildSubscriptionTierThree,
		discordgo.MessageTypeUserPremiumGuildSubscription:
	// TODO: booster present
	case
		discordgo.MessageTypeDefault, discordgo.MessageTypeReply:
		// TODO: mute player
	}

	if strings.HasPrefix(m.Content, "!") {

		var vars []string
		var content string
		inputSplit := strings.Split(m.Content, " ")
		for idx := range inputSplit {
			if idx == 0 {
				content = inputSplit[idx]
			} else {
				vars = append(vars, inputSplit[idx])
			}
		}
		switch content {
		case "!help_zxc":
			d.printWelcome(m.Author.ID, m.GuildID)
		case "!login_zxc":
			if d.isAdmin(m.ChannelID, m.Author.ID) {
				d.printLogin(m.ChannelID)
			}
		case "!update_zxc":
			d.checkUpdate()
		}
	}
	return
}

func (d *Discord) onUserChanged(s *discordgo.Session, i *discordgo.GuildMemberUpdate) {
	return
}

func (d *Discord) onUserConnected(s *discordgo.Session, i *discordgo.GuildMemberAdd) {
	d.logger.Infof("user %s connected to server %s", i.User.Username, i.GuildID)
	d.printWelcome(i.User.ID, i.GuildID)
	return
}

func (d *Discord) onUserDisconnected(s *discordgo.Session, i *discordgo.GuildMemberRemove) {
	d.logger.Infof("user %s disconnected from server %s", i.User.Username, i.GuildID)
	d.deleteUser(i.User.ID)
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
	case discordgo.InteractionMessageComponent:
		if h, ok := d.components()[i.MessageComponentData().CustomID]; ok {
			h(s, i)
		}
	}
}

func (d *Discord) refreshAll() {
	ticker := time.NewTicker(60 * time.Second)
	quit := make(chan struct{})
	for {
		select {
		case <-ticker.C:
			go d.checkUpdate()
			go d.listenQueue()
			go d.updateStats()
		case <-quit:
			ticker.Stop()
			return
		}
	}
}
