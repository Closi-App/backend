package main

import (
	"context"
	"flag"
	"github.com/Closi-App/backend/cmd/app/wire"
	"github.com/Closi-App/backend/internal/config"
)

func main() {
	configFilePath := flag.String("config", "./config.yml", "config file path, eg: -config ./configs/local.yml")
	flag.Parse()
	cfg := config.Load(*configFilePath)

	app, cleanup, err := wire.NewWire(cfg)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	app.Run(context.Background())
}
