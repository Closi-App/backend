package service

import (
	"context"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/repository"
	"go.mongodb.org/mongo-driver/v2/bson"
	"time"
)

type AnswerService interface {
	Create(ctx context.Context, input AnswerCreateInput) (bson.ObjectID, error)
	GetAll(ctx context.Context, filter ...domain.AnswerGetAllFilter) ([]domain.Answer, error)
	GetByID(ctx context.Context, id bson.ObjectID) (domain.Answer, error)
	Update(ctx context.Context, id, userID bson.ObjectID, input domain.AnswerUpdateInput) error
	Delete(ctx context.Context, id, userID bson.ObjectID) error

	AddLike(ctx context.Context, id bson.ObjectID) error
	RemoveLike(ctx context.Context, id bson.ObjectID) error
	Verify(ctx context.Context, id bson.ObjectID) error
}

type answerService struct {
	*Service
	repository      repository.AnswerRepository
	questionService QuestionService
	userService     UserService
}

func NewAnswerService(
	service *Service,
	repository repository.AnswerRepository,
	questionService QuestionService,
	userService UserService,
) AnswerService {
	return &answerService{
		Service:         service,
		repository:      repository,
		questionService: questionService,
		userService:     userService,
	}
}

type AnswerCreateInput struct {
	Text       string
	QuestionID bson.ObjectID
	UserID     bson.ObjectID
}

func (s *answerService) Create(ctx context.Context, input AnswerCreateInput) (bson.ObjectID, error) {
	id := bson.NewObjectID()

	if err := s.repository.Create(ctx, domain.Answer{
		ID:         id,
		Text:       input.Text,
		Likes:      0,
		QuestionID: input.QuestionID,
		UserID:     input.UserID,
		IsVerified: false,
		CreatedAt:  time.Now(),
		UpdatedAt:  time.Now(),
	}); err != nil {
		return bson.ObjectID{}, err
	}

	question, err := s.questionService.GetByID(ctx, input.QuestionID)
	if err != nil {
		return bson.ObjectID{}, err
	}

	if err = s.userService.AdjustPoints(ctx, input.UserID, int(question.Points)); err != nil {
		return bson.ObjectID{}, err
	}

	if err = s.userService.AdjustPoints(ctx, question.UserID, -1*int(question.Points)); err != nil {
		return bson.ObjectID{}, err
	}

	return id, nil
}

func (s *answerService) GetAll(ctx context.Context, filter ...domain.AnswerGetAllFilter) ([]domain.Answer, error) {
	return s.repository.GetAll(ctx, filter...)
}

func (s *answerService) GetByID(ctx context.Context, id bson.ObjectID) (domain.Answer, error) {
	return s.repository.GetByID(ctx, id)
}

func (s *answerService) Update(ctx context.Context, id, userID bson.ObjectID, input domain.AnswerUpdateInput) error {
	return s.repository.Update(ctx, id, userID, input)
}

func (s *answerService) Delete(ctx context.Context, id, userID bson.ObjectID) error {
	return s.repository.Delete(ctx, id, userID)
}

func (s *answerService) AddLike(ctx context.Context, id bson.ObjectID) error {
	return s.repository.AddLike(ctx, id)
}

func (s *answerService) RemoveLike(ctx context.Context, id bson.ObjectID) error {
	return s.repository.RemoveLike(ctx, id)
}

func (s *answerService) Verify(ctx context.Context, id bson.ObjectID) error {
	return s.repository.Verify(ctx, id)
}
