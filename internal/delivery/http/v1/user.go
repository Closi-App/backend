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
		users.Post("/refresh", h.userRefresh)
		users.Get("/:id", h.userGetByID)

		auth := users.Group("", h.authMiddleware)
		{
			auth.Get("/", h.userGet)
			auth.Put("/", h.userUpdate)
			auth.Delete("/", h.userDelete)
		}
	}
}

type userSignUpRequest struct {
	Name         string          `json:"name"`
	Username     string          `json:"username"`
	Email        string          `json:"email"`
	Password     string          `json:"password"`
	Location     domain.Location `json:"location"`
	Language     domain.Language `json:"language"`
	ReferrerCode string          `json:"referrer_code"`
}

type userSignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Sign up
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
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		Location:     req.Location,
		Language:     req.Language,
		ReferrerCode: req.ReferrerCode,
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

// @Summary		Sign in
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

type userRefreshRequest struct {
	Token string `json:"token"`
}

type userRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Refresh tokens
// @Description	Refresh the user's access and refresh tokens
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userRefreshRequest	body		userRefreshRequest	true	"Request"
// @Success		200					{object}	userRefreshResponse
// @Failure		400,500				{object}	errorResponse
// @Router			/users/refresh [post]
func (h *Handler) userRefresh(ctx *fiber.Ctx) error {
	var req userRefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.userService.RefreshTokens(ctx.Context(), req.Token)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, userRefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		Get current user
// @Description	Retrieve the currently authenticated user's information
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200			{object}	domain.User
// @Failure		400,401,500	{object}	errorResponse
// @Router			/users [get]
func (h *Handler) userGet(ctx *fiber.Ctx) error {
	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	user, err := h.userService.GetByID(ctx.Context(), ctxUser.ID)
	if err != nil {
		if errors.Is(err, domain.ErrUserNotFound) {
			return h.newResponse(ctx, fiber.StatusBadRequest, err)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, user)
}

// @Summary		Get user by ID
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
	Name      string              `json:"name"`
	Username  string              `json:"username"`
	Email     string              `json:"email"`
	Password  string              `json:"password"`
	AvatarURL string              `json:"avatar_url"`
	Settings  domain.UserSettings `json:"settings"`
}

// @Summary		Update current user
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
	var req userUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err := h.userService.Update(ctx.Context(), ctxUser.ID, service.UserUpdateInput{
		Name:      req.Name,
		Username:  req.Username,
		Email:     req.Email,
		Password:  req.Password,
		AvatarURL: req.AvatarURL,
		Settings:  req.Settings,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Delete current user
// @Description	Delete the authenticated user's account
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Success		200		{string}	string	"OK"
// @Failure		401,500	{object}	errorResponse
// @Router			/users [delete]
func (h *Handler) userDelete(ctx *fiber.Ctx) error {
	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err := h.userService.Delete(ctx.Context(), ctxUser.ID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
