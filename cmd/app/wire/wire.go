//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/logger"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewUserRepository,
)

func NewWire(*config.Config) (*repository.Repository, func(), error) {
	panic(wire.Build(
		logger.NewZerolog,
		repositorySet,
	))
}
