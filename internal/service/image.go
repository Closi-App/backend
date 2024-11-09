package service

import (
	"context"
	"github.com/Closi-App/backend/internal/repository"
)

type ImageService interface {
	Upload(ctx context.Context, fileBytes []byte) (url string, err error)
}

type imageService struct {
	*Service
	repository repository.ImageRepository
}

func NewImageService(service *Service, repository repository.ImageRepository) ImageService {
	return &imageService{
		Service:    service,
		repository: repository,
	}
}

func (s *imageService) Upload(ctx context.Context, fileBytes []byte) (string, error) {
	return s.repository.Upload(ctx, fileBytes)
}
