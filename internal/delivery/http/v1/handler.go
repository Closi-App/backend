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
	router.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	h.initUserRoutes(router)
}
