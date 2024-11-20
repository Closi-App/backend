package v1

import (
	"github.com/Closi-App/backend/internal/service"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/Closi-App/backend/pkg/localizer"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
	"golang.org/x/text/language"
)

type Handler struct {
	log                   *logger.Logger
	localizer             *localizer.Localizer
	countryService        service.CountryService
	imageService          service.ImageService
	tagService            service.TagService
	userService           service.UserService
	questionService       service.QuestionService
	answerService         service.AnswerService
	tokensManager         auth.TokensManager
	appSupportedLanguages []language.Tag
}

func NewHandler(
	log *logger.Logger,
	localizer *localizer.Localizer,
	countryService service.CountryService,
	imageService service.ImageService,
	tagService service.TagService,
	userService service.UserService,
	questionService service.QuestionService,
	answerService service.AnswerService,
	tokensManager auth.TokensManager,
	appSupportedLanguages []language.Tag,
) *Handler {
	return &Handler{
		log:                   log,
		localizer:             localizer,
		countryService:        countryService,
		imageService:          imageService,
		tagService:            tagService,
		userService:           userService,
		questionService:       questionService,
		answerService:         answerService,
		tokensManager:         tokensManager,
		appSupportedLanguages: appSupportedLanguages,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	v1 := router.Group("/v1", h.localizerMiddleware)
	{
		h.initCountryRoutes(v1)
		h.initImageRoutes(v1)
		h.initTagRoutes(v1)
		h.initUserRoutes(v1)
		h.initQuestionRoutes(v1)
		h.initAnswerRoutes(v1)
	}
}
