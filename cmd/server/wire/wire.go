//go:build wireinject
// +build wireinject

package wire

import (
	"github.com/Closi-App/backend/internal/app"
	"github.com/Closi-App/backend/internal/delivery/http"
	"github.com/Closi-App/backend/internal/delivery/http/v1"
	"github.com/Closi-App/backend/internal/repository"
	"github.com/Closi-App/backend/internal/service"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/Closi-App/backend/pkg/database/mongo"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

var pkgSet = wire.NewSet(
	mongo.NewMongo,
	auth.NewTokensManager,
	auth.NewPasswordHasher,
)

var repositorySet = wire.NewSet(
	repository.NewRepository,
	repository.NewUserRepository,
	repository.NewQuestionRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewUserService,
	service.NewQuestionService,
)

var deliverySet = wire.NewSet(
	v1.NewHandler,
	http.NewServer,
)

func newApp(cfg *viper.Viper, log *logger.Logger, httpServer *http.Server) *app.App {
	return app.NewApp(cfg, log, httpServer)
}

func NewWire(*viper.Viper, *logger.Logger) (*app.App, func(), error) {
	panic(wire.Build(
		pkgSet,
		repositorySet,
		serviceSet,
		deliverySet,
		newApp,
	))
}
