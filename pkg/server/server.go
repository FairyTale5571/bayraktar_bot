package server

import (
	"github.com/fairytale5571/bayraktar_bot/pkg/cache"
	"net/http"

	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/links"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
	"github.com/gin-gonic/gin"
)

type Router struct {
	router *gin.Engine
	cfg    *models.Config
	bot    *discord.Discord
	logger *logger.Wrapper
	db     *database.DB
	rdb    *redis.Redis
	cache  *cache.Config
}

func New(cfg *models.Config, bot *discord.Discord, db *database.DB, rdb *redis.Redis) *Router {
	r := Router{
		bot:    bot,
		cfg:    cfg,
		db:     db,
		rdb:    rdb,
		router: gin.Default(),
		logger: logger.New("server"),
		cache:  cache.SetupCache(rdb),
	}
	return &r
}

func (r *Router) Start() {
	r.logger.Info("gin opened")
	r.router.LoadHTMLGlob("webApp/static/*.html")

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

	r.router.GET("/site", r.site)

	apiGroup := r.router.Group("/api")
	{
		apiGame := apiGroup.Group("/game")
		{
			apiGame.GET("/adminRules", r.adminRules)
			apiGame.GET("/news", r.news)

			apiGame.GET("/economy", r.economy)
			apiGame.GET("/gov", r.government)
			apiTop := apiGame.Group("/top")
			{
				apiTop.GET("/players", r.topPlayer)
				apiTop.GET("/gangs", r.topGang)
				apiTop.GET("/wanted", r.wanted)
			}
		}

		apiDiscord := apiGroup.Group("/discord")
		{
			apiDiscord.POST("/mailing/:guild", r.mailingUsers) // not use, bot will be banned
			apiDiscord.POST("/direct/:userid", r.sendDirect)
			apiDiscord.POST("/channel/:guild/:channel", r.sendToChannel)

			apiDiscord.GET("/guilds", nil)
			apiDiscord.GET("/channels/:guild", nil)
			apiDiscord.GET("/members/:guild", nil)
			apiDiscord.GET("/roles/:guild", nil)
		}
	}
}

func (r *Router) plugin(c *gin.Context) {
	c.Redirect(http.StatusTemporaryRedirect, r.cfg.URL+"/assets/files/task_force_radio.ts3_plugin")
}

func (r *Router) site(c *gin.Context) {
	c.HTML(http.StatusOK, "/assets/views/login_button.html", nil)
}

func (r *Router) topPlayer(c *gin.Context) {

}

func (r *Router) topGang(c *gin.Context) {

}

func (r *Router) wanted(c *gin.Context) {

}

func (r *Router) adminRules(c *gin.Context) {
	c.HTML(http.StatusOK, "admin_rules.html", nil)
}
