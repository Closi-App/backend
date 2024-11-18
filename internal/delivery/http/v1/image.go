package v1

import (
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"io"
)

func (h *Handler) initImageRoutes(router fiber.Router) {
	images := router.Group("/images", h.userAuthMiddleware)
	{
		images.Post("/", h.imageUpload)
	}
}

// TODO: rename this handler "files"

// @Summary		Upload
// @Description	Upload image
// @Security		UserAuth
// @Tags			images
// @Accept			multipart/form-data
// @Produce		json
// @Param			image		formData	file	true	"Image"
// @Success		201			{object}	successResponse
// @Failure		400,401,500	{object}	errorResponse
// @Router			/images [post]
func (h *Handler) imageUpload(ctx *fiber.Ctx) error {
	f, err := ctx.FormFile("image")
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	file, err := f.Open()
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}
	defer file.Close()

	fileBytes, err := io.ReadAll(file)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	url, err := h.imageService.Upload(ctx.Context(), fileBytes)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, domain.ErrInternalServerError)
	}

	return h.newResponse(ctx, fiber.StatusCreated, urlResponse{url})
}
