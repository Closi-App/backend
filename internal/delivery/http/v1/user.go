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
		users.Delete("/", h.authMiddleware, h.userDelete)
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

// @Summary		Sign Up
// @Description	Create a new user account
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userSignUpRequest	body		userSignUpRequest	true	"Request"
// @Success		200					{object}	userSignUpResponse
// @Failure		400,409,500			{object}	errorResponse
// @Router			/users/sign-up [post]
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
			return h.newResponse(ctx, fiber.StatusConflict, err)
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

// @Summary		Sign In
// @Description	Authenticate a user and retrieve tokens
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userSignInRequest	body		userSignInRequest	true	"Request"
// @Success		200					{object}	userSignInResponse
// @Failure		400,500				{object}	errorResponse
// @Router			/users/sign-in [post]
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

// @Summary		Get Current User
// @Description	Retrieve the currently authenticated user's information
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200			{object}	domain.User
// @Failure		400,401,500	{object}	errorResponse
// @Router			/users [get]
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

// @Summary		Get User by ID
// @Description	Retrieve a user's information by their ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"User ID"
// @Success		200			{object}	domain.User
// @Failure		400,404,500	{object}	errorResponse
// @Router			/users/{id} [get]
func (h *Handler) userGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	user, err := h.userService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusNotFound, err)
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

// @Summary		Update Current User
// @Description	Update the authenticated user's information
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userUpdateRequest	body		userUpdateRequest	true	"Request"
// @Success		200					{string}	string				"OK"
// @Failure		400,401,500			{object}	errorResponse
// @Router			/users [put]
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

// @Summary		Delete Current User
// @Description	Delete the authenticated user's account
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200		{string}	string	"OK"
// @Failure		401,500	{object}	errorResponse
// @Router			/users [delete]
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
