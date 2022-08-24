package app

import (
	"github.com/caarlos0/env/v6"
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/discord"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
	"github.com/fairytale5571/bayraktar_bot/pkg/server"
	"github.com/fairytale5571/bayraktar_bot/pkg/storage/redis"
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

	rdb, err := redis.New(cfg.RedisUri)
	if err != nil {
		log.Fatalf("cant create redis client: %v", err)
		return nil, err
	}
	log.Info("redis started")

	ds, err := discord.New(cfg, db, rdb)
	if err != nil {
		log.Fatalf("error start discord: %v", err)
		return nil, err
	}

	srv := server.New(cfg, ds, db, rdb)
	if err != nil {
		log.Errorf("error start server: %v", err)
		return nil, err
	}

	go srv.Start()
	go ds.Start()
	log.Info("application started")
	return &App{
		Discord: ds,
		DB:      db,
		Config:  cfg,
		Logger:  log,
		Server:  srv,
	}, nil
}
