package main

import (
	"log"
	"os"
	"os/signal"

	"github.com/fairytale5571/bayraktar_bot/pkg/app"
)

func main() {
	app, err := app.New()
	if err != nil {
		log.Fatalf("start application failed: %v", err)
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	sig := <-stop
	app.Logger.Infof("Close: received %v", sig.String())

	err = app.DB.Close()
	if err != nil {
		app.Logger.Errorf("Close: error close database: %v", err)
		return
	}
	app.Discord.Stop()

	log.Fatalf("Graceful shutdown\n************************************************************************\n\n")
}
