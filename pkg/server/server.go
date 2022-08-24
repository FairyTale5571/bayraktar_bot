package server

import (
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
	"net/http"
	"strings"

	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/gin-gonic/gin"
	"golang.org/x/oauth2"
)

type Router struct {
	router   *gin.Engine
	cfg      *models.Config
	bot      *discord.Discord
	logger   *logger.LoggerWrapper
	db       *database.DB
	rdb      *redis.Redis
	settings AuthConfig
}

type AuthConfig struct {
	DiscordConfig oauth2.Config
}

func New(cfg *models.Config, bot *discord.Discord, db *database.DB, rdb *redis.Redis) *Router {
	r := Router{
		bot:    bot,
		cfg:    cfg,
		db:     db,
		rdb:    rdb,
		router: gin.Default(),
		logger: logger.New("server"),
	}

	r.settings = AuthConfig{
		DiscordConfig: models.DiscordOauth,
	}
	return &r
}

func (r *Router) Start() {
	r.logger.Info("gin opened")
	r.router.Static("/assets/", "webApp/assets/")
	r.mainRouter()
	err := r.router.Run(":" + r.cfg.PORT)
	if err != nil {
		r.logger.Errorf("cant open gin engine: %v", err)
		return
	}
}

func (r *Router) Stop() {
	// TODO: implement
}

func (r *Router) mainRouter() {
	r.router.GET("/", func(c *gin.Context) {
		c.Redirect(http.StatusPermanentRedirect, links.UrlSite)
	})
	r.router.GET("/auth/steam", r.steam)
	r.router.GET("/redirect/:to", r.redirect)
	r.router.GET("/plugin", r.plugin)

	apiGroup := r.router.Group("/api")
	{
		apiGroup.GET("/economy", r.economy)

	}
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

func (r *Router) plugin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, r.cfg.URL+"/assets/files/task_force_radio.ts3_plugin")
}
