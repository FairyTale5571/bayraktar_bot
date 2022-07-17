package server

import (
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
	settings AuthConfig
}

type AuthConfig struct {
	DiscordConfig oauth2.Config
}

func New(cfg *models.Config, bot *discord.Discord) *Router {
	r := Router{
		bot:    bot,
		cfg:    cfg,
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

}
