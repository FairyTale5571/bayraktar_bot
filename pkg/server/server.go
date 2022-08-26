package server

import (
	"net/http"

	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
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
	authGroup := r.router.Group("/auth")
	{
		authGroup.GET("/steam", r.steam)
		authGroup.GET("/discord", r.discord)
	}
	r.router.GET("/redirect/:to", r.redirect)
	r.router.GET("/plugin", r.plugin)

	apiGroup := r.router.Group("/api")
	{
		apiGroup.GET("/economy", r.economy)
		apiGroup.POST("/mailing/:guild", r.mailingUsers)
		apiGroup.POST("/direct/:userid", r.sendDirect)
		apiGroup.POST("/channel/:guild/:channel", r.sendToChannel)
	}
}

func (r *Router) plugin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, r.cfg.URL+"/assets/files/task_force_radio.ts3_plugin")
}
