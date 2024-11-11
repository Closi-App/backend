package v1

import (
	"github.com/Closi-App/backend/internal/service"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	log             *logger.Logger
	countryService  service.CountryService
	imageService    service.ImageService
	tagService      service.TagService
	userService     service.UserService
	questionService service.QuestionService
	answerService   service.AnswerService
	tokensManager   auth.TokensManager
}

func NewHandler(
	log *logger.Logger,
	countryService service.CountryService,
	imageService service.ImageService,
	tagService service.TagService,
	userService service.UserService,
	questionService service.QuestionService,
	answerService service.AnswerService,
	tokensManager auth.TokensManager,
) *Handler {
	return &Handler{
		log:             log,
		countryService:  countryService,
		imageService:    imageService,
		tagService:      tagService,
		userService:     userService,
		questionService: questionService,
		answerService:   answerService,
		tokensManager:   tokensManager,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	{
		h.initCountryRoutes(v1)
		h.initImageRoutes(v1)
		h.initTagRoutes(v1)
		h.initUserRoutes(v1)
		h.initQuestionRoutes(v1)
		h.initAnswerRoutes(v1)
	}
}
