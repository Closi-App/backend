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
	"github.com/Closi-App/backend/pkg/database/redis"
	"github.com/Closi-App/backend/pkg/imgbb"
	"github.com/Closi-App/backend/pkg/localizer"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/Closi-App/backend/pkg/smtp"
	"github.com/google/wire"
	"github.com/spf13/viper"
	"golang.org/x/text/language"
)

var pkgSet = wire.NewSet(
	localizer.NewLocalizer,
	logger.NewLogger,
	mongo.NewMongo,
	redis.NewRedis,
	imgbb.NewImgbb,
	smtp.NewSMTPSender,
	auth.NewTokensManager,
	auth.NewPasswordHasher,
)

var repositorySet = wire.NewSet(
	repository.NewRepository,
	repository.NewCountryRepository,
	repository.NewImageRepository,
	repository.NewTagRepository,
	repository.NewUserRepository,
	repository.NewQuestionRepository,
	repository.NewAnswerRepository,
)

var serviceSet = wire.NewSet(
	service.NewService,
	service.NewCountryService,
	service.NewImageService,
	service.NewEmailService,
	service.NewTagService,
	service.NewUserService,
	service.NewQuestionService,
	service.NewAnswerService,
)

var deliverySet = wire.NewSet(
	v1.NewHandler,
	http.NewServer,
)

func newApp(cfg *viper.Viper, log *logger.Logger, httpServer *http.Server) *app.App {
	return app.NewApp(cfg, log, httpServer)
}

func NewWire(*viper.Viper, []language.Tag) (*app.App, func(), error) {
	panic(wire.Build(
		pkgSet,
		repositorySet,
		serviceSet,
		deliverySet,
		newApp,
	))
}
