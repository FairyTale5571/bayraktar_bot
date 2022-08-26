package server

import (
	"net/http"
	"strings"

	"github.com/gin-gonic/gin"
)

func (r *Router) discord(c *gin.Context) {
	body, err := c.GetRawData()
	if err != nil {
		c.JSON(http.StatusInternalServerError, gin.H{"error": err.Error()})
		return
	}
	r.logger.Infof("discord auth: %v", string(body))
}

func (r *Router) steam(c *gin.Context) {
	r.logger.Info("steam auth")

	state := c.Request.URL.Query().Get("state")
	guild := c.Request.URL.Query().Get("guild")
	steamId := c.Request.URL.Query().Get("openid.claimed_id")
	if steamId == "" {
		r.logger.Error("steam auth failed")
		c.JSON(http.StatusForbidden, gin.H{
			"message": "steam auth failed",
		})
		return
	}
	steamId = strings.TrimLeft(steamId, "https://steamcommunity.com/openid/id/")
	r.logger.Infof("steam auth: %v | %v", state, steamId)
	r.bot.RegisterUser(guild, state, steamId)
	c.String(http.StatusOK, "Проверьте сообщение от бота в личных сообщениях")
}
