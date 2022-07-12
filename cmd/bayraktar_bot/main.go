package main

import (
	"github.com/fairytale5571/bayraktar_bot/pkg/app"
	"log"
	"os"
	"os/signal"
)

func main() {
	if err := app.New(); err != nil {
		log.Fatalf("start application failed: %v", err)
	}
	stop := make(chan os.Signal)
	signal.Notify(stop, os.Interrupt)
	<-stop
	log.Fatalf("Graceful shutdown\n************************************************************************\n\n")
}
