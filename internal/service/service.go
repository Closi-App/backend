package service

import "github.com/Closi-App/backend/pkg/logger"

type Service struct {
	log *logger.Logger
}

func NewService(log *logger.Logger) *Service {
	return &Service{
		log: log,
	}
}
