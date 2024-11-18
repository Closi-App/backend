package v1

import (
	"errors"
	"fmt"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
)

const (
	successStatus responseStatus = "success"
	errorStatus   responseStatus = "error"
)

type responseStatus string

type response struct {
	Status     responseStatus `json:"status"`
	StatusCode int            `json:"status_code"`
	RequestID  string         `json:"request_id"`
}

type successResponse struct {
	response
	Data interface{} `json:"data"`
}

type errorResponse struct {
	response
	Error domain.Error `json:"error"`
}

type messageResponse struct {
	Message string `json:"message"`
}

type idResponse struct {
	ID string `json:"id"`
}

type urlResponse struct {
	URL string `json:"url"`
}

func (h *Handler) newSuccessResponse(ctx *fiber.Ctx, res response, data interface{}) error {
	return ctx.Status(res.StatusCode).JSON(successResponse{
		response: res,
		Data:     data,
	})
}

func (h *Handler) newErrorResponse(ctx *fiber.Ctx, res response, err error) error {
	var appErr *domain.Error
	if errors.As(err, &appErr) {
		localizer := h.getLocalizerFromCtx(ctx)

		return ctx.Status(res.StatusCode).JSON(errorResponse{
			response: res,
			Error: domain.Error{
				Code:    fmt.Sprintf("errors.%s", appErr.Code),
				Message: localizer.Translate(appErr.Code),
			},
		})
	} else {
		h.log.Error().
			Str("requestID", h.getRequestIDFromCtx(ctx)).
			Err(err).
			Msg("Internal server error")
		return h.newErrorResponse(ctx, res, domain.ErrInternalServerError)
	}
}

func (h *Handler) newResponse(ctx *fiber.Ctx, statusCode int, v ...interface{}) error {
	res := response{
		StatusCode: statusCode,
		RequestID:  h.getRequestIDFromCtx(ctx),
	}

	if statusCode < 300 {
		res.Status = successStatus
	} else {
		res.Status = errorStatus
	}

	if len(v) > 0 {
		value := v[0]

		switch value.(type) {
		case error:
			return h.newErrorResponse(ctx, res, value.(error))
		case string:
			return h.newSuccessResponse(ctx, res, messageResponse{value.(string)})
		default:
			return h.newSuccessResponse(ctx, res, value)
		}
	}

	return ctx.Status(statusCode).JSON(res)
}
