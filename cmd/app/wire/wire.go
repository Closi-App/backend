//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Closi-App/backend/internal/config"
	"github.com/Closi-App/backend/internal/delivery/http"
	"github.com/Closi-App/backend/internal/delivery/http/v1"
	"github.com/Closi-App/backend/internal/logger"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/Closi-App/backend/internal/service"
	"github.com/Closi-App/backend/pkg/app"
	"github.com/google/wire"
)

var repositorySet = wire.NewSet(
	repository.NewDB,
	repository.NewRepository,
	repository.NewUserRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
)

var deliverySet = wire.NewSet(
	v1.NewHandler,
	http.NewServer,
)

func newApp(httpServer *http.Server) *app.App {
	return app.NewApp(httpServer)
}

func NewWire(*config.Config) (*app.App, func(), error) {
	panic(wire.Build(
		logger.NewZerolog,
		repositorySet,
		serviceSet,
		deliverySet,
		newApp,
	))
}
