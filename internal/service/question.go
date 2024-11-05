package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type QuestionService interface {
	Create(ctx context.Context, input QuestionCreateInput) (string, error)
	Get(ctx context.Context) ([]domain.Question, error)
	GetByLocation(ctx context.Context, location domain.Location) ([]domain.Question, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input QuestionUpdateInput) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error
}

type questionService struct {
	*Service
	repository repository.QuestionRepository
}

func NewQuestionService(service *Service, repository repository.QuestionRepository) QuestionService {
	return &questionService{
		Service:    service,
		repository: repository,
	}
}

type QuestionCreateInput struct {
	Title       string
	Description string
	Attachments []string
	Points      uint
	Location    domain.Location
	UserID      bson.ObjectID
}

func (s *questionService) Create(ctx context.Context, input QuestionCreateInput) (string, error) {
	id := bson.NewObjectID()

	if err := s.repository.Create(ctx, domain.Question{
		ID:          id,
		Title:       input.Title,
		Description: input.Description,
		Attachments: input.Attachments,
		Points:      input.Points,
		Location:    input.Location,
		UserID:      input.UserID,
		CreatedAt:   time.Now(),
		UpdatedAt:   time.Now(),
	}); err != nil {
		return "", err
	}

	return id.Hex(), nil
}

func (s *questionService) Get(ctx context.Context) ([]domain.Question, error) {
	return s.repository.Get(ctx)
}

func (s *questionService) GetByLocation(ctx context.Context, location domain.Location) ([]domain.Question, error) {
	return s.repository.GetByLocation(ctx, location)
}

func (s *questionService) GetByID(ctx context.Context, id bson.ObjectID) (domain.Question, error) {
	return s.repository.GetByID(ctx, id)
}

type QuestionUpdateInput struct {
	Title       string
	Description string
	Attachments []string
	Points      uint
	Location    domain.Location
}

func (s *questionService) Update(ctx context.Context, id, userID bson.ObjectID, input QuestionUpdateInput) error {
	return s.repository.Update(ctx, id, userID, repository.QuestionUpdateInput{
		Title:       input.Title,
		Description: input.Description,
		Attachments: input.Attachments,
		Points:      input.Points,
		Location:    &input.Location,
	})
}

func (s *questionService) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	return s.repository.Delete(ctx, id, userID)
}
