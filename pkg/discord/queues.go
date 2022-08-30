package discord

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"

	"github.com/fairytale5571/bayraktar_bot/pkg/links"
)

func (d *Discord) getGroupsRole(id int) (string, string) {
	if id == -1 {
		return "", ""
	}
	rows, err := d.db.Query("SELECT id, ds_role_leader, ds_role_member_id FROM groups WHERE id = ?", id)
	defer rows.Close()

	if err != nil {
		d.logger.Errorf("getGroupsRole(): Error: %v", err.Error())
		return "", ""
	}
	var groupId uint8
	var dsRoleLeader, dsRoleMember string
	for rows.Next() {
		if err := rows.Scan(&groupId, &dsRoleLeader, &dsRoleMember); err != nil {
			d.logger.Errorf("getGroupsRole(): Error: %v", err.Error())
			continue
		}
		return dsRoleLeader, dsRoleMember
	}
	return "", ""
}

func (d *Discord) isLeaderGroup(id int, steamID string) bool {
	var creator, leader string
	rows, err := d.db.Query("SELECT creator, leader FROM groups WHERE id = ?", id)
	defer rows.Close()
	if err != nil {
		d.logger.Errorf("isLeaderGroup(): Error: %v", err.Error())
		return false
	}
	for rows.Next() {
		if err := rows.Scan(&creator, &leader); err != nil {
			d.logger.Errorf("isLeaderGroup(): Error: %v", err.Error())
			continue
		}
		if creator == steamID || leader == steamID {
			return true
		}
	}
	return false
}

func (d *Discord) listenQueue() {
	type player struct {
		ds         string
		steam      string
		donorLevel int
		groupId    int
	}

	_getQueue := func() []*player {
		var queue []*player

		rows, err := d.db.Query(`select discord_users.discord_uid, discord_users.uid, p.donorlevel, p.group_id from discord_users 
    		inner join players p on discord_users.uid = p.playerid		
    		inner join discord_queue dq on discord_users.uid = dq.uid order by dq.id asc`)
		defer rows.Close()
		if err != nil {
			d.logger.Errorf("Error getting queue: %s", err.Error())
			return nil
		}
		for rows.Next() {
			var t player
			if err := rows.Scan(&t.ds, &t.steam, &t.donorLevel, &t.groupId); err != nil {
				d.logger.Errorf("Error getting queue: %s", err.Error())
				return nil
			}
			queue = append(queue, &t)
		}
		_, _ = d.db.Exec("TRUNCATE TABLE discord_queue")
		return queue
	}

	_rename := func(v *player) {
		var name, fullName sql.NullString
		err := d.db.QueryRow("select name, CONCAT(first_name, ' ', last_name) from players where playerid = ?", v.steam).Scan(&name, &fullName)
		if err != nil {
			d.logger.Errorf("Error getting player name (%s): %s", v.steam, err.Error())
			return
		}
		if fullName.Valid {
			err = d.ds.GuildMemberNickname(d.cfg.GuildID, v.ds, fullName.String)
			if err != nil {
				d.logger.Errorf("Error setting player name (%s): %s", v.steam, err.Error())
				return
			}
			return
		}
		err = d.ds.GuildMemberNickname(d.cfg.GuildID, v.ds, name.String)
		if err != nil {
			d.logger.Errorf("Error setting player name (%s): %s", v.steam, err.Error())
			return
		}
	}
	_reRole := func(v *player) {
		_, err := d.ds.GuildMember(d.cfg.GuildID, v.ds)
		if err != nil {
			d.deleteUserFromVerified(v.ds)
			return
		}
		if !d.haveRole(d.cfg.GuildID, d.cfg.RegRoleID, v.ds) {
			err := d.setRole(d.cfg.GuildID, d.cfg.RegRoleID, v.ds)
			if err != nil {
				d.logger.Errorf("Error setting role (%s): %s", v.steam, err.Error())
				return
			}
		}
		if v.donorLevel > 0 {
			if !d.haveRole(d.cfg.GuildID, d.cfg.VipRole, v.ds) {
				err := d.setRole(d.cfg.GuildID, d.cfg.VipRole, v.ds)
				if err != nil {
					d.logger.Errorf("Error setting role (%s): %s", v.steam, err.Error())
					return
				}
			}
		} else {
			if d.haveRole(d.cfg.GuildID, d.cfg.VipRole, v.ds) {
				d.removeRole(d.cfg.GuildID, d.cfg.VipRole, v.ds)
			}
		}
		groupRoles := d.getGroupRoles()
		for _, gr := range groupRoles {
			if d.haveRole(d.cfg.GuildID, gr, v.ds) {
				d.removeRole(d.cfg.GuildID, gr, v.ds)
			}
		}

		if v.groupId > 0 {
			lead, member := d.getGroupsRole(v.groupId)
			if d.isLeaderGroup(v.groupId, v.steam) {
				err := d.setRole(d.cfg.GuildID, lead, v.ds)
				if err != nil {
					d.logger.Errorf("Error setting role (%s): %s", v.steam, err.Error())
					return
				}
			} else {
				err := d.setRole(d.cfg.GuildID, member, v.ds)
				if err != nil {
					d.logger.Errorf("Error setting role (%s): %s", v.steam, err.Error())
					return
				}
			}
		}
	}

	for _, v := range _getQueue() {
		_rename(v)
		_reRole(v)
	}
}

func (d *Discord) getLkApi() (*gov, error) {
	var players gov
	var client http.Client
	resp, err := client.Get(links.UrlLk + "/api/gov")
	defer resp.Body.Close() // nolint: not needed

	if err != nil {
		log.Println(err)
		return nil, err
	}

	if resp.StatusCode == http.StatusOK {
		bodyBytes, err := io.ReadAll(resp.Body)
		if err != nil {
			log.Println(err)
			return nil, err
		}
		if !json.Valid(bodyBytes) {
			fmt.Println("json is invalid")
			return nil, err
		}
		err = json.Unmarshal(bodyBytes, &players)
		if err != nil {
			fmt.Println(err)
			return nil, err
		}
	}
	return &players, nil
}
