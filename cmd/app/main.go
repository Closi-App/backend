package main

import (
	"github.com/Closi-App/backend/cmd/app/wire"
	"github.com/Closi-App/backend/internal/config"
)

func main() {
	cfg := config.Load("./config.yml")

	log, cleanup, err := wire.NewWire(cfg)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	log.Info("OK")
}
