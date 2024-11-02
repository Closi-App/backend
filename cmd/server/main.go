package main

import (
	"context"
	"flag"
	"github.com/Closi-App/backend/cmd/server/wire"
	"github.com/Closi-App/backend/pkg/config"
	"github.com/Closi-App/backend/pkg/logger"
)

func main() {
	cfgFilePath := flag.String("config", "./config.yml", "config file path, eg: -config ./configs/local.yml")
	flag.Parse()
	cfg := config.NewConfig(*cfgFilePath)

	log := logger.NewLogger(cfg)

	app, cleanup, err := wire.NewWire(cfg, log)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	app.Run(context.Background())
}
