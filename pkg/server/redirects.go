package server

import (
	"net/http"

	"github.com/fairytale5571/bayraktar_bot/pkg/links"
	"github.com/gin-gonic/gin"
)

/*
	redirect.GET("/discord", r.discordRedirect)
	redirect.GET("/forum", r.forumRedirect)
	redirect.GET("/site", r.siteRedirect)
	redirect.GET("/lk", r.lkRedirect)
	redirect.GET("/mod", r.modRedirect)
	redirect.GET("/plugin", r.pluginRedirect)
*/

var linksMap = map[string]string{
	"discord":   links.UrlDiscord,
	"forum":     links.UrlForum,
	"site":      links.UrlSite,
	"lk":        links.UrlLk,
	"mod":       links.UrlMod,
	"game":      links.UrlGame,
	"teamspeak": links.UrlTeamspeak,
}

func (r *Router) redirect(c *gin.Context) {
	to := c.Param("to")
	if v, ok := linksMap[to]; ok {
		c.Redirect(http.StatusTemporaryRedirect, v)
		return
	}
	c.Redirect(http.StatusTemporaryRedirect, links.UrlSite)
}
