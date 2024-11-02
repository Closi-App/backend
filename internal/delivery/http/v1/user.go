package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
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
	var req userSignUpRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.userService.SignUp(ctx.Context(), service.UserSignUpInput{
		Name:     req.Name,
		Username: req.Username,
		Email:    req.Email,
		Password: req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserAlreadyExists) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, userSignUpResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
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
	var req userSignInRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.userService.SignIn(ctx.Context(), service.UserSignInInput{
		UsernameOrEmail: req.UsernameOrEmail,
		Password:        req.Password,
	})
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, userSignInResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

func (h *Handler) userGet(ctx *fiber.Ctx) error {
	id, err := h.getUserID(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	user, err := h.userService.GetByID(ctx.Context(), id)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, user)
}

func (h *Handler) userGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	user, err := h.userService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, user)
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
	id, err := h.getUserID(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	var req userUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	if err := h.userService.Update(ctx.Context(), id, service.UserUpdateInput{
		Name:                    req.Name,
		Username:                req.Username,
		Email:                   req.Email,
		Password:                req.Password,
		AvatarURL:               req.AvatarURL,
		NotificationPreferences: req.NotificationPreferences,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

func (h *Handler) userDelete(ctx *fiber.Ctx) error {
	id, err := h.getUserID(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err := h.userService.Delete(ctx.Context(), id); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
