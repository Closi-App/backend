package v1

import (
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) authMiddleware(ctx *fiber.Ctx) error {
	return ctx.Next() // TODO
}

func (h *Handler) getUserID(ctx *fiber.Ctx) bson.ObjectID {
	id, err := bson.ObjectIDFromHex("") // TODO
	if err != nil {
		return bson.NilObjectID
	}

	return id
}
