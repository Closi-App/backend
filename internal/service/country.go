package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type CountryService interface {
	Create(ctx context.Context, input CountryCreateInput) (bson.ObjectID, error)
	Get(ctx context.Context) ([]domain.Country, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Country, error)
}

type countryService struct {
	*Service
	repository repository.CountryRepository
}

func NewCountryService(service *Service, repository repository.CountryRepository) CountryService {
	return &countryService{
		Service:    service,
		repository: repository,
	}
}

type CountryCreateInput struct {
	Name map[domain.Language]string
}

func (s *countryService) Create(ctx context.Context, input CountryCreateInput) (bson.ObjectID, error) {
	id := bson.NewObjectID()

	if err := s.repository.Create(ctx, domain.Country{
		ID:   id,
		Name: input.Name,
	}); err != nil {
		return bson.ObjectID{}, err
	}

	return id, nil
}

func (s *countryService) Get(ctx context.Context) ([]domain.Country, error) {
	return s.repository.Get(ctx)
}

func (s *countryService) GetByID(ctx context.Context, id bson.ObjectID) (domain.Country, error) {
	return s.repository.GetByID(ctx, id)
}
