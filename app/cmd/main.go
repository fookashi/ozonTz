package main

import (
	"context"
	"os"
	"os/signal"
	"syscall"

	"app/internal/app"
	"app/internal/config"
)

func main() {
	cfg := config.MustLoadConfig()

	application := app.NewApp(context.Background(), cfg)

	go func() {
		application.HttpApp.Run()
	}()

	// graceful stop
	stop := make(chan os.Signal, 1)
	signal.Notify(stop, syscall.SIGINT, syscall.SIGTERM)

	<-stop

	application.HttpApp.Stop()
}
