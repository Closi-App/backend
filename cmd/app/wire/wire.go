//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/logger"
	"github.com/google/wire"
)

func NewWire(*config.Config) (logger.Logger, func(), error) {
	panic(wire.Build(
		logger.NewZerolog,
	))
}
