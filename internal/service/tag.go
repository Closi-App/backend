package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
)

type TagService interface {
	Create(ctx context.Context, input TagCreateInput) (bson.ObjectID, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Tag, error)
	GetAll(ctx context.Context) ([]domain.Tag, error)
	GetAllByCountryID(ctx context.Context, countryID bson.ObjectID) ([]domain.Tag, error)
	Delete(ctx context.Context, id bson.ObjectID) error
}

type tagService struct {
	*Service
	repository repository.TagRepository
}

func NewTagService(service *Service, repository repository.TagRepository) TagService {
	return &tagService{
		Service:    service,
		repository: repository,
	}
}

type TagCreateInput struct {
	Name      string
	CountryID bson.ObjectID
}

func (s *tagService) Create(ctx context.Context, input TagCreateInput) (bson.ObjectID, error) {
	id := bson.NewObjectID()

	if err := s.repository.Create(ctx, domain.Tag{
		ID:        id,
		Name:      input.Name,
		CountryID: input.CountryID,
	}); err != nil {
		return bson.ObjectID{}, err
	}

	return id, nil
}

func (s *tagService) GetByID(ctx context.Context, id bson.ObjectID) (domain.Tag, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *tagService) GetAll(ctx context.Context) ([]domain.Tag, error) {
	return s.repository.GetAll(ctx)
}

func (s *tagService) GetAllByCountryID(ctx context.Context, countryID bson.ObjectID) ([]domain.Tag, error) {
	return s.repository.GetAllByCountryID(ctx, countryID)
}

func (s *tagService) Delete(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Delete(ctx, id)
}
