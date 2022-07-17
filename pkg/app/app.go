package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/fairytale5571/bayraktar_bot/pkg/server"
)

type App struct {
	Discord *discord.Discord
	DB      *database.DB
	Config  *models.Config
	Logger  *logger.LoggerWrapper
	Server  *server.Router
}

func New() (*App, error) {
	log := logger.New("application")

	cfg := &models.Config{}
	if err := env.Parse(cfg); err != nil {
		log.Errorf("error parse config: %v", err)
		return nil, err
	}

	db, err := database.New(cfg.MysqlUri)
	if err != nil {
		log.Errorf("error start database: %v", err)
		return nil, err
	}

	ds, err := discord.New(cfg, db)
	if err != nil {
		log.Fatalf("error start discord: %v", err)
		return nil, err
	}

	server := server.New(cfg, ds)
	if err != nil {
		log.Errorf("error start server: %v", err)
		return nil, err
	}

	go server.Start()
	go ds.Start()
	log.Info("application started")
	return &App{
		Discord: ds,
		DB:      db,
		Config:  cfg,
		Logger:  log,
		Server:  server,
	}, nil
}
