package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/gin-gonic/gin"
)

const (
	headerPass = "X-Pass"
)

// @Summary Get game
// @Description get current economy
// @Tags game
// @Success 	200 			{object} 			models.Economy
// @Failure 	500 			{object} 			models.Error
// @Router /api/game/economy [get]
func (r *Router) economy(c *gin.Context) {
	rows, err := r.db.Query("SELECT * FROM economy")
	defer rows.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	type Economy struct {
		LastUpdate       time.Time
		Resource         string
		Localize         string
		Influenced       string
		Price            int
		MaxPrice         int
		RandomMax        int
		RandomMin        int
		MinPrice         int
		DownPricePerItem float64
		RandomDownPrice  bool
		Illegal          bool
	}
	var economies []Economy
	for rows.Next() {
		var e Economy
		err := rows.Scan(&e.Resource, &e.Localize, &e.Price, &e.MaxPrice, &e.MinPrice, &e.DownPricePerItem, &e.RandomDownPrice, &e.RandomMax, &e.RandomMin, &e.Illegal, &e.Influenced, &e.LastUpdate)
		if err != nil {
			c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
			return
		}
		economies = append(economies, e)
	}
	c.JSON(http.StatusOK, economies)
}

// @Summary Discord send to channel
// @Description send message to channel
// @Tags discord
// @Accept json
// @Produce json
// @Param guild path string true "guild id"
// @Param channel path string true "channel id"
// @Param X-Pass header string true "password"
// @Success 	200
// @Failure 	500 			{object} 			models.Error
// @Router /api/discord/send/channel/{guild}/{channel} [post]
func (r *Router) sendToChannel(c *gin.Context) {
	if c.GetHeader(headerPass) != r.cfg.PostPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "pass not valid"})
		return
	}
	guild := c.Param("guild")
	if guild == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guild not found"})
		return
	}
	channel := c.Param("channel")
	if channel == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guild not found"})
		return
	}

	body_, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var embeds models.Embeds
	err = json.Unmarshal(body_, &embeds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	go r.bot.SendToChannel(guild, channel, embeds)
}

// @Summary Discord send to user
// @Description send message to user
// @Tags discord
// @Accept json
// @Produce json
// @Param user path string true "user id"
// @Param X-Pass header string true "password"
// @Param body body models.Embeds true "embeds"
// @Success 	200
// @Failure 	500 			{object} 			models.Error
// @Router /api/discord/direct/{userid} [post]
func (r *Router) sendDirect(c *gin.Context) {
	if c.GetHeader(headerPass) != r.cfg.PostPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "pass not valid"})
		return
	}
	user := c.Param("userid")
	if user == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "discord user not found"})
		return
	}

	body_, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	var embeds models.Embeds
	err = json.Unmarshal(body_, &embeds)
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
	})
	go r.bot.SendDirect(user, embeds)
}

func (r *Router) government(c *gin.Context) {
	var govResponse models.Government
	govs, err := r.db.Query(`
		select count(uid) c, "police" t from players where side = "cop" and connected = 1
    	union select count(uid) c, "all" t from players where connected = 1
    union select count(uid) c, "ems" t from players where side = "med" and connected = 1
    union select count(uid) c, "rev" t from players where side = "reb" and connected = 1
    union select count(uid) c, "civ" t from players where side = "civ" and connected = 1;
	`)
	defer govs.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"point": "govs",
		})
		return
	}

	groupList, err := r.db.Query(`
		select id, premial_var from groups`)
	defer groupList.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"point": "groups",
		})
		return
	}
	plrActiveList, err := r.db.Query(`
		select group_id from players where connected = 1 and group_id > 0`)
	defer plrActiveList.Close()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{
			"error": err.Error(),
			"point": "players",
		})
		return
	}
	c.JSON(http.StatusOK, gin.H{
		"message": "ok",
		"gov":     govResponse,
	})
}
