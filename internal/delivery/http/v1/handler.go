package v1

import (
	"github.com/Closi-App/backend/internal/logger"
	"github.com/Closi-App/backend/internal/service"
	"github.com/gofiber/fiber/v2"
)

type Handler struct {
	log         logger.Logger
	userService service.UserService
}

func NewHandler(log logger.Logger, userService service.UserService) *Handler {
	return &Handler{
		log:         log,
		userService: userService,
	}
}

func (h *Handler) InitRoutes(router fiber.Router) {
	router.Get("/ping", func(ctx *fiber.Ctx) error {
		return ctx.SendString("pong")
	})

	h.initUserRoutes(router)
}
