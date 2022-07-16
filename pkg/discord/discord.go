package discord

import (
	"github.com/bwmarrin/discordgo"
	"github.com/fairytale5571/bayraktar_bot/pkg/database"
	"github.com/fairytale5571/bayraktar_bot/pkg/logger"
	"github.com/fairytale5571/bayraktar_bot/pkg/models"
)

type Discord struct {
	cfg    *models.Config
	logger *logger.LoggerWrapper
	ds     *discordgo.Session
	db     *database.DB
}

func New(cfg *models.Config, db *database.DB) (*Discord, error) {
	res := &Discord{
		cfg:    cfg,
		db:     db,
		logger: logger.New("discord"),
	}
	s, err := discordgo.New("Bot " + res.cfg.DiscordToken)
	if err != nil {
		res.logger.Fatalf("cant create discord session: %v", err)
		return nil, err
	}
	res.ds = s
	return res, nil
}

func (d *Discord) Start() {
	d.ds.Identify.Intents = discordgo.MakeIntent(discordgo.IntentsAll)
	d.setupHandlers()

	err := d.ds.Open()
	if err != nil {
		d.logger.Fatalf("cant open discord session: %v", err)
		return
	}
	d.logger.Info("discord started")
}

func (d *Discord) setupHandlers() {
	d.ds.AddHandler(d.onBotUp)
	d.ds.AddHandler(d.onMessageCreate)
	d.ds.AddHandler(d.onUserChanged)
	d.ds.AddHandler(d.onCommandsCall)
	d.ds.AddHandler(d.onUserConnected)
	d.ds.AddHandler(d.onUserDisconnected)
	d.ds.AddHandler(d.onGuildCreate)
}

func (d *Discord) Stop() {
	err := d.ds.Close()
	if err != nil {
		d.logger.Errorf("cant close discord session: %v", err)
		return
	}
}
