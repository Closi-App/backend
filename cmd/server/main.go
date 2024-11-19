package main

import (
	"context"
	"flag"
	"github.com/Closi-App/backend/cmd/server/wire"
	"github.com/Closi-App/backend/pkg/config"
)

//	@title			Closi API
//	@version		1.0
//	@description	REST API for Closi App
//	@host			127.0.0.1:8080
//	@BasePath		/api/v1/

// @securityDefinitions.apikey	UserAuth
// @in							header
// @name						Authorization
func main() {
	cfgFilePath := flag.String("config", "./config.yml", "config file path, eg: -config ./configs/local.yml")
	flag.Parse()
	cfg := config.NewConfig(*cfgFilePath)

	app, cleanup, err := wire.NewWire(cfg)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	app.Run(context.Background())
}
