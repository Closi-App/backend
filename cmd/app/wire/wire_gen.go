// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

package wire

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/logger"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/google/wire"
)

// Injectors from wire.go:

func NewWire(configConfig *config.Config) (*repository.Repository, func(), error) {
	loggerLogger := logger.NewZerolog(configConfig)
	client := repository.NewDB(configConfig, loggerLogger)
	repositoryRepository := repository.New(loggerLogger, client)
	return repositoryRepository, func() {
	}, nil
}

// wire.go:

var repositorySet = wire.NewSet(repository.NewDB, repository.New)
