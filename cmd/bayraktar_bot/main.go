package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/fairytale5571/bayraktar_bot/pkg/app"
)

// @title Bayraktar Bot
// @version 1.0
// @description API for Bayraktar bot

// @host localhost:3200
// @BasePath /
// @Schemes http https

func main() {
	a, err := app.New()
	if err != nil {
		log.Fatalf("start application failed: %v", err)
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	a.Logger.Infof("Close: received %v", sig.String())

	err = a.DB.Close()
	if err != nil {
		a.Logger.Errorf("Close: error close database: %v", err)
		return
	}
	a.Discord.Stop()
	a.Server.Stop()
	log.Fatalf("Graceful shutdown\n************************************************************************\n\n")
}
