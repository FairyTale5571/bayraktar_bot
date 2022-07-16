package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
)

type App struct {
	Discord *discord.Discord
	DB      *database.DB
	Config  *models.Config
	Logger  *logger.LoggerWrapper
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

	go ds.Start()
	log.Info("application started")
	return &App{
		Discord: ds,
		DB:      db,
		Config:  cfg,
		Logger:  log,
	}, nil
}
