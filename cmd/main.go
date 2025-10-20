package main

import (
	"auth_service/internal/bootstrap"
	"auth_service/internal/config"
	"context"
	"os"
	"os/signal"

	"github.com/joho/godotenv"
	"github.com/sethvargo/go-envconfig"
)

func main() {
	_ = godotenv.Load(".env")
	var cfg config.Config

	err := envconfig.ProcessWith(context.TODO(), &envconfig.Config{Target: &cfg, Lookuper: envconfig.OsLookuper()})
	if err != nil {
		panic(err)
	}

	quitSignal := make(chan os.Signal, 1)
	signal.Notify(quitSignal, os.Interrupt)

	app := bootstrap.New(cfg)

	ctx, cancel := context.WithCancel(context.Background())
	go func() {
		<-quitSignal
		cancel()
	}()

	app.Run(ctx)
}
