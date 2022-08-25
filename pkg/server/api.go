package server

import (
	"encoding/json"
	"net/http"
	"time"

	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/gin-gonic/gin"
)

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

const (
	headerPass  = "X-Pass"
	discordUser = "X-Discord-User"
)

func (r *Router) mailingUsers(c *gin.Context) {
	guild := c.Param("guild")
	if guild == "" {
		c.JSON(http.StatusBadRequest, gin.H{"error": "guild not found"})
		return
	}
	if c.GetHeader(headerPass) != r.cfg.PostPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "pass not valid"})
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
	go r.bot.SendMassive(guild, embeds)
}

func (r *Router) sendDirect(c *gin.Context) {
	if c.GetHeader(headerPass) != r.cfg.PostPassword {
		c.JSON(http.StatusUnauthorized, gin.H{"error": "pass not valid"})
		return
	}
	user := c.GetHeader(discordUser)
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
