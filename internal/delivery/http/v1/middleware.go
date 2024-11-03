package v1

import (
	"errors"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
	"strings"
)

const (
	authorizationHeader = "Authorization"

	userCtxKey = "userID"
)

func (h *Handler) authMiddleware(ctx *fiber.Ctx) error {
	header := ctx.Get(authorizationHeader)
	if header == "" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	accessToken := headerParts[1]

	userID, err := h.tokensManager.Parse(accessToken)
	if err != nil {
		return ctx.SendStatus(fiber.StatusUnauthorized)
	}

	ctx.Locals(userCtxKey, userID)

	return ctx.Next()
}

func (h *Handler) getUserID(ctx *fiber.Ctx) (bson.ObjectID, error) {
	id, ok := ctx.Locals(userCtxKey).(string)
	if !ok {
		return bson.ObjectID{}, errors.New("error getting user id from context")
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return bson.ObjectID{}, errors.New("invalid object id")
	}

	return objectID, nil
}
