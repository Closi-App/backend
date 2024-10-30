package v1

import (
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

func (h *Handler) initUserRoutes(router fiber.Router) {
	users := router.Group("/users")
	{
		users.Post("/sign-up", h.userSignUp)
		users.Post("/sign-in", h.userSignIn)
		users.Get("/", h.authMiddleware, h.userGet)
		users.Get("/:id", h.userGetByID)
		users.Put("/", h.authMiddleware, h.userUpdate)
	}
}

type userSignUpRequest struct {
	Name     string `json:"name"`
	Username string `json:"username"`
	Email    string `json:"email"`
	Password string `json:"password"`
}

type userSignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) userSignUp(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

type userSignInRequest struct {
	UsernameOrEmail string `json:"username_or_email"`
	Password        string `json:"password"`
}

type userSignInResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

func (h *Handler) userSignIn(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) userGet(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) userGetByID(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

type userUpdateRequest struct {
	Name                    string                         `json:"name"`
	Username                string                         `json:"username"`
	Email                   string                         `json:"email"`
	Password                string                         `json:"password"`
	AvatarURL               string                         `json:"avatar_url"`
	NotificationPreferences domain.NotificationPreferences `json:"notification_preferences"`
}

func (h *Handler) userUpdate(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}

func (h *Handler) userDelete(ctx *fiber.Ctx) error {
	return ctx.SendStatus(fiber.StatusOK)
}
