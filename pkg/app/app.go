package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
)

func New() error {
	log := logger.New("application")

	cfg := models.Config{}
	if err := env.Parse(&cfg); err != nil {
		log.Errorf("error parse config: %v", err)
		return err
	}

	db, err := database.New(cfg.MysqlUri)
	if err != nil {
		log.Errorf("error start database: %v", err)
		return err
	}

	ds, err := discord.New(&cfg, db)
	if err != nil {
		log.Fatalf("error start discord: %v", err)
		return err
	}

	go ds.Start()
	log.Info("application started")
	return nil
}
