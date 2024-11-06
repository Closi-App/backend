package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
)

const (
	authorizationHeader = "Authorization"
	userCtxKey          = "user"
)

func (h *Handler) parseUserIDFromAuthHeader(ctx *fiber.Ctx) (string, error) {
	header := ctx.Get(authorizationHeader)
	if header == "" {
		return "", errors.New("empty authorization header")
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return "", errors.New("invalid authorization header")
	}

	accessToken := headerParts[1]

	userID, err := h.tokensManager.Parse(accessToken)
	if err != nil {
		return "", errors.New("invalid access token")
	}

	return userID, nil
}

func (h *Handler) authUserMiddleware(ctx *fiber.Ctx) error {
	id, err := h.parseUserIDFromAuthHeader(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	user, err := h.userService.GetByID(ctx.Context(), objectID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	ctx.Locals(userCtxKey, user)

	return ctx.Next()
}

func (h *Handler) getUserFromCtx(ctx *fiber.Ctx) (domain.User, error) {
	user, ok := ctx.Locals(userCtxKey).(domain.User)
	if !ok {
		return domain.User{}, errors.New("error getting user from context")
	}

	return user, nil
}
