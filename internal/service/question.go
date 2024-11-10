package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type QuestionService interface {
	Create(ctx context.Context, input QuestionCreateInput) (bson.ObjectID, error)
	GetAll(ctx context.Context, filter ...domain.QuestionGetAllFilter) ([]domain.Question, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input domain.QuestionUpdateInput) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
}

type questionService struct {
	*Service
	repository repository.QuestionRepository
	tagService TagService
}

func NewQuestionService(service *Service, repository repository.QuestionRepository, tagService TagService) QuestionService {
	return &questionService{
		Service:    service,
		repository: repository,
		tagService: tagService,
	}
}

type QuestionCreateInput struct {
	Title          string
	Description    string
	AttachmentsURL []string
	Tags           []string
	Points         uint
	CountryID      bson.ObjectID
	UserID         bson.ObjectID
}

func (s *questionService) Create(ctx context.Context, input QuestionCreateInput) (bson.ObjectID, error) {
	var tags []bson.ObjectID

	for _, tagName := range input.Tags {
		id, err := s.tagService.Create(ctx, TagCreateInput{
			Name:      tagName,
			CountryID: input.CountryID,
		})
		if err != nil {
			return bson.ObjectID{}, err
		}

		tags = append(tags, id)
	}

	id := bson.NewObjectID()

	if err := s.repository.Create(ctx, domain.Question{
		ID:             id,
		Title:          input.Title,
		Description:    input.Description,
		AttachmentsURL: input.AttachmentsURL,
		Tags:           tags,
		Points:         input.Points,
		CountryID:      input.CountryID,
		UserID:         input.UserID,
		CreatedAt:      time.Now(),
		UpdatedAt:      time.Now(),
	}); err != nil {
		return bson.ObjectID{}, err
	}

	return id, nil
}

func (s *questionService) GetAll(ctx context.Context, filter ...domain.QuestionGetAllFilter) ([]domain.Question, error) {
	return s.repository.GetAll(ctx, filter...)
}

func (s *questionService) GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *questionService) Update(ctx context.Context, id, userID bson.ObjectID, input domain.QuestionUpdateInput) error {
	return s.repository.Update(ctx, id, userID, input)
}

func (s *questionService) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	return s.repository.Delete(ctx, id, userID)
}
