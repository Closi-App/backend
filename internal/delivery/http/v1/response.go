package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

type idResponse struct {
	ID string `json:"id"`
}

type urlResponse struct {
	URL string `json:"url"`
}

type errorResponse struct {
	Code    string `json:"code"`
	Message string `json:"message"`
}

func (h *Handler) newErrorResponse(ctx *fiber.Ctx, statusCode int, err error) error {
	var appErr *domain.Error
	if errors.As(err, &appErr) {
		return ctx.Status(statusCode).JSON(errorResponse{
			Code:    appErr.Code,
			Message: appErr.Message,
		})
	} else {
		h.log.Error().
			Err(err).
			Msg("Internal server error")
		return h.newErrorResponse(ctx, statusCode, domain.ErrInternalServerError)
	}
}

func (h *Handler) newResponse(ctx *fiber.Ctx, statusCode int, data interface{}) error {
	switch data.(type) {
	case error:
		return h.newErrorResponse(ctx, statusCode, data.(error))
	case nil:
		return ctx.SendStatus(statusCode)
	default:
		return ctx.Status(statusCode).JSON(data)
	}
}
