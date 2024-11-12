// Code generated by Wire. DO NOT EDIT.

//go:generate go run -mod=mod github.com/google/wire/cmd/wire
//go:build !wireinject
// +build !wireinject

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
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/Closi-App/backend/pkg/smtp"
	"github.com/google/wire"
	"github.com/spf13/viper"
)

// Injectors from wire.go:

func NewWire(viperViper *viper.Viper, loggerLogger *logger.Logger) (*app.App, func(), error) {
	serviceService := service.NewService(loggerLogger)
	database := mongo.NewMongo(viperViper)
	client := redis.NewRedis(viperViper)
	imgbbClient := imgbb.NewImgbb(viperViper)
	repositoryRepository := repository.NewRepository(loggerLogger, database, client, imgbbClient)
	countryRepository := repository.NewCountryRepository(repositoryRepository)
	countryService := service.NewCountryService(serviceService, countryRepository)
	imageRepository := repository.NewImageRepository(repositoryRepository)
	imageService := service.NewImageService(serviceService, imageRepository)
	tagRepository := repository.NewTagRepository(repositoryRepository)
	tagService := service.NewTagService(serviceService, tagRepository)
	userRepository := repository.NewUserRepository(repositoryRepository)
	sender := smtp.NewSMTPSender(viperViper)
	emailService := service.NewEmailService(serviceService, sender)
	passwordHasher := auth.NewPasswordHasher(viperViper)
	tokensManager := auth.NewTokensManager(viperViper)
	userService := service.NewUserService(serviceService, viperViper, userRepository, emailService, passwordHasher, tokensManager)
	questionRepository := repository.NewQuestionRepository(repositoryRepository)
	questionService := service.NewQuestionService(serviceService, questionRepository, tagService)
	answerRepository := repository.NewAnswerRepository(repositoryRepository)
	answerService := service.NewAnswerService(serviceService, answerRepository, questionService, userService)
	handler := v1.NewHandler(loggerLogger, countryService, imageService, tagService, userService, questionService, answerService, tokensManager)
	server := http.NewServer(viperViper, loggerLogger, handler)
	appApp := newApp(viperViper, loggerLogger, server)
	return appApp, func() {
	}, nil
}

// wire.go:

var pkgSet = wire.NewSet(mongo.NewMongo, redis.NewRedis, imgbb.NewImgbb, smtp.NewSMTPSender, auth.NewTokensManager, auth.NewPasswordHasher)

var repositorySet = wire.NewSet(repository.NewRepository, repository.NewCountryRepository, repository.NewImageRepository, repository.NewTagRepository, repository.NewUserRepository, repository.NewQuestionRepository, repository.NewAnswerRepository)

var serviceSet = wire.NewSet(service.NewService, service.NewCountryService, service.NewImageService, service.NewEmailService, service.NewTagService, service.NewUserService, service.NewQuestionService, service.NewAnswerService)

var deliverySet = wire.NewSet(v1.NewHandler, http.NewServer)

func newApp(cfg *viper.Viper, log *logger.Logger, httpServer *http.Server) *app.App {
	return app.NewApp(cfg, log, httpServer)
}
