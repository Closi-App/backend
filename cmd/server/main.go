package main

import (
	"context"
	"flag"
	"github.com/Closi-App/backend/cmd/server/wire"
	"github.com/Closi-App/backend/pkg/config"
	"github.com/Closi-App/backend/pkg/localizer"
	"github.com/Closi-App/backend/pkg/logger"
	"golang.org/x/text/language"
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

	// TODO: think about place of localizer initializer
	l := localizer.NewLocalizer([]string{
		"./locales/en.json",
		"./locales/uk.json",
		"./locales/de.json",
		"./locales/pl.json",
		"./locales/ru.json",
	}, language.English)

	log := logger.NewLogger(cfg)

	app, cleanup, err := wire.NewWire(cfg, l, log)
	defer cleanup()
	if err != nil {
		panic(err)
	}

	app.Run(context.Background())
}
