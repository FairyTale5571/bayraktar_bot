package discord

import (
	"database/sql"
	"encoding/json"
	"fmt"
	"io"
	"log"
	"net/http"
)

const guild = "719969719871995958"

func (d *Discord) listenQueue() {

	type player struct {
		ds    string
		steam string
	}

	_getQueue := func() []*player {
		var queue []*player

		rows, err := d.db.Query("select discord_users.discord_uid, discord_users.uid from discord_users inner join discord_queue dq on discord_users.uid = dq.uid")
		//defer rows.Close()
		if err != nil {
			d.logger.Errorf("Error getting queue: %s", err.Error())
			return nil
		}
		for rows.Next() {
			var t player
			if err := rows.Scan(&t.ds, &t.steam); err != nil {
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
			err = d.ds.GuildMemberNickname(guild, v.ds, fullName.String)
			if err != nil {
				d.logger.Errorf("Error setting player name (%s): %s", v.steam, err.Error())
				return
			}
			return
		}
		err = d.ds.GuildMemberNickname(guild, v.ds, name.String)
		if err != nil {
			d.logger.Errorf("Error setting player name (%s): %s", v.steam, err.Error())
			return
		}

	}
	_reRole := func(v *player) {

	}

	for _, v := range _getQueue() {
		_rename(v)
		_reRole(v)
	}

}

func (d *Discord) updateStats() {
	gov, _ := d.getLkApi()
	_, err := d.ds.ChannelEdit("953006266303787019", fmt.Sprintf("║ 🌐 На сервере: %d", gov.Gov.Info.All))
	if err != nil {
		d.logger.Errorf("Error updating channel: %s", err.Error())
		return
	}
}

func (d *Discord) getLkApi() (*gov, error) {
	var players gov
	var client http.Client
	resp, err := client.Get("https://lk.rimasrp.life/api/gov")
	if err != nil {
		log.Println(err)
		return nil, err
	}
	defer resp.Body.Close()

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
