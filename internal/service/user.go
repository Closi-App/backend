package service

import "github.com/Closi-App/backend/internal/repository"

type UserService interface {
}

type userService struct {
	*Service
	repository repository.UserRepository
}

func NewUserService(service *Service, repository repository.UserRepository) UserService {
	return &userService{
		Service:    service,
		repository: repository,
	}
}
