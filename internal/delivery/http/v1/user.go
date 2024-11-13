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
		users.Post("/:id/confirm", h.userConfirm)

		auth := users.Group("", h.authUserMiddleware)
		{
			auth.Get("/", h.userGet)
			auth.Put("/", h.userUpdate)
			auth.Delete("/", h.userDelete)

			favorites := auth.Group("/favorites")
			{
				favorites.Post("/:questionID", h.userAddFavorite)
				favorites.Delete("/:questionID", h.userRemoveFavorite)
			}
		}

		// TODO: admins handlers
	}
}

// TODO: fields validation

type userSignUpRequest struct {
	Name         string `json:"name"`
	Username     string `json:"username"`
	Email        string `json:"email"`
	Password     string `json:"password"`
	CountryID    string `json:"country_id"`
	Language     string `json:"language"`
	ReferrerCode string `json:"referrer_code"`
}

type userSignUpResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Sign up
// @Description	Sign up
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

	countryObjectID, err := bson.ObjectIDFromHex(req.CountryID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	language := domain.ParseLanguage(req.Language)

	tokens, err := h.userService.SignUp(ctx.Context(), service.UserSignUpInput{
		Name:         req.Name,
		Username:     req.Username,
		Email:        req.Email,
		Password:     req.Password,
		CountryID:    countryObjectID,
		Language:     language,
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
// @Description	Sign in
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

// @Summary		Get
// @Description	Get auth user
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

// @Summary		Get by ID
// @Description	Get user by ID
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
	Name               *string
	Username           *string
	Email              *string
	Password           *string
	AvatarURL          *string
	CountryID          *string
	Language           *string
	Appearance         *string
	EmailNotifications *bool
}

// @Summary		Update
// @Description	Update auth user
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userUpdateRequest	body		userUpdateRequest	true	"Request"
// @Success		200					{string}	string				"OK"
// @Failure		400,401,500			{object}	errorResponse
// @Router			/users [put]
func (h *Handler) userUpdate(ctx *fiber.Ctx) error {
	var err error

	var req userUpdateRequest
	if err = ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	var (
		countryObjectID bson.ObjectID
		language        domain.Language
		appearance      domain.Appearance
	)

	if req.CountryID != nil {
		countryObjectID, err = bson.ObjectIDFromHex(*req.CountryID)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}
	}
	if req.Language != nil {
		language = domain.ParseLanguage(*req.Language)
	}
	if req.Appearance != nil {
		appearance = domain.ParseAppearance(*req.Appearance)
	}

	if err = h.userService.Update(ctx.Context(), ctxUser.ID, domain.UserUpdateInput{
		Name:               req.Name,
		Username:           req.Username,
		Email:              req.Email,
		Password:           req.Password,
		AvatarURL:          req.AvatarURL,
		CountryID:          &countryObjectID,
		Language:           &language,
		Appearance:         &appearance,
		EmailNotifications: req.EmailNotifications,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Delete
// @Description	Delete auth user
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

	if err = h.userService.Delete(ctx.Context(), ctxUser.ID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

type userRefreshRequest struct {
	Token string `json:"token"`
}

type userRefreshResponse struct {
	AccessToken  string `json:"access_token"`
	RefreshToken string `json:"refresh_token"`
}

// @Summary		Refresh tokens
// @Description	Refresh user's tokens
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			userRefreshRequest	body		userRefreshRequest	true	"Request"
// @Success		200					{object}	userRefreshResponse
// @Failure		400,401,500			{object}	errorResponse
// @Router			/users/refresh [post]
func (h *Handler) userRefresh(ctx *fiber.Ctx) error {
	var req userRefreshRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	tokens, err := h.userService.RefreshTokens(ctx.Context(), req.Token)
	if err != nil {
		if errors.Is(err, domain.ErrUnauthorized) {
			return h.newResponse(ctx, fiber.StatusUnauthorized, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, userRefreshResponse{
		AccessToken:  tokens.AccessToken,
		RefreshToken: tokens.RefreshToken,
	})
}

// @Summary		Confirm by ID
// @Description	Confirm user by ID
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			id		path		string	true	"User ID"
// @Success		200		{string}	string	"OK"
// @Failure		400,500	{object}	errorResponse
// @Router			/users/{id}/confirm [get]
func (h *Handler) userConfirm(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	if err = h.userService.Confirm(ctx.Context(), objectID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Add favorite
// @Description	Add favorite for auth user
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			questionID	path		string	true	"Question ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/users/favorites/{questionID} [post]
func (h *Handler) userAddFavorite(ctx *fiber.Ctx) error {
	questionID := ctx.Params("questionID")
	questionObjectID, err := bson.ObjectIDFromHex(questionID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err = h.userService.AddFavorite(ctx.Context(), ctxUser.ID, questionObjectID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Remove favorite
// @Description	Remove favorite from auth user
// @Security		UserAuth
// @Tags			users
// @Accept			json
// @Produce		json
// @Param			questionID	path		string	true	"Question ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/users/favorites/{questionID} [delete]
func (h *Handler) userRemoveFavorite(ctx *fiber.Ctx) error {
	questionID := ctx.Params("questionID")
	questionObjectID, err := bson.ObjectIDFromHex(questionID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err = h.userService.RemoveFavorite(ctx.Context(), ctxUser.ID, questionObjectID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
