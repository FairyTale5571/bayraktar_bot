package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/errorUtils"
)

func (d *Discord) setRole(guildID, roleID, userID string) error {
	_, err := d.ds.GuildMember(guildID, userID)
	if err != nil {
		return err
	}
	err = d.ds.GuildMemberRoleAdd(guildID, userID, roleID)
	if err != nil {
		return err
	}
	return nil
}

func (d *Discord) haveRole(guildID, roleID, userID string) bool {
	member, err := d.ds.GuildMember(guildID, userID)
	if err != nil {
		d.logger.Infof("User: %v will be deleted", userID)
		d.deleteUserFromVerified(userID)
		return false
	}
	for _, role := range member.Roles {
		if role == roleID {
			return true
		}
	}
	return false
}

func (d *Discord) removeRole(guildID, roleID, userID string) {
	_, err := d.ds.GuildMember(guildID, userID)
	if err != nil {
		d.logger.Infof("User: %v will be deleted", userID)
		d.deleteUserFromVerified(userID)
		return
	}
	err = d.ds.GuildMemberRoleRemove(guildID, userID, roleID)
	if err != nil {
		return
	}
	return
}

func (d *Discord) findRoleById(guildID, roleID string) (*discordgo.Role, error) {
	roles, err := d.ds.GuildRoles(guildID)
	if err != nil {
		return nil, err
	}
	for _, role := range roles {
		if role.ID == roleID {
			return role, nil
		}
	}
	return nil, errorUtils.ErrRoleNotFound
}
