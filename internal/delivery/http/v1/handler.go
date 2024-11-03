package v1

import (
	"github.com/Closi-App/backend/internal/service"
	"github.com/Closi-App/backend/pkg/auth"
	"github.com/Closi-App/backend/pkg/logger"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	log           *logger.Logger
	userService   service.UserService
	tokensManager auth.TokensManager
}

func NewHandler(log *logger.Logger, userService service.UserService, tokensManager auth.TokensManager) *Handler {
	return &Handler{
		log:           log,
		userService:   userService,
		tokensManager: tokensManager,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	v1 := router.Group("/v1")
	{
		h.initUserRoutes(v1)
	}
}
