package main

import (
	"flag"
	"github.com/Closi-App/backend/cmd/server/wire"
	"github.com/Closi-App/backend/internal/config"
)

func main() {
	configFilePath := flag.String("config", "./config.yml", "config file path, eg: -config ./configs/local.yml")
	flag.Parse()
	cfg := config.Load(*configFilePath)

	_, cleanup, err := wire.NewWire(cfg)
	defer cleanup()
	if err != nil {
		panic(err)
	}
}
