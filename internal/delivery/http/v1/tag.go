package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) initTagRoutes(router fiber.Router) {
	tags := router.Group("/tags")
	{
		tags.Get("/:id", h.tagGetByID)
		tags.Get("/", h.tagGetAll)

		auth := tags.Group("", h.authUserMiddleware)
		{
			auth.Post("/", h.tagCreate)
			tags.Get("/country/:countryID", h.tagGetAllByCountryID)
		}

		// TODO: delete function for admins
	}
}

type tagCreateRequest struct {
	Name      string `json:"name"`
	CountryID string `json:"country_id"`
}

// @Summary		Create
// @Description	Create new tag
// @Security		UserAuth
// @Tags			tags
// @Accept			json
// @Produce		json
// @Param			tagCreateRequest	body		tagCreateRequest	true	"Request"
// @Success		201					{object}	idResponse
// @Failure		400,401,500			{object}	errorResponse
// @Router			/tags [post]
func (h *Handler) tagCreate(ctx *fiber.Ctx) error {
	var req tagCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	id, err := h.tagService.Create(ctx.Context(), service.TagCreateInput{
		Name:      req.Name,
		CountryID: ctxUser.Settings.CountryID,
	})
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, domain.ErrBadRequest)
	}

	return h.newResponse(ctx, fiber.StatusCreated, idResponse{id.Hex()})
}

// @Summary		Get by ID
// @Description	Get tag by ID
// @Tags			tags
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Tag ID"
// @Success		200			{object}	domain.Tag
// @Failure		400,404,500	{object}	errorResponse
// @Router			/tags/{id} [get]
func (h *Handler) tagGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	tag, err := h.tagService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrTagNotFound) {
			return h.newResponse(ctx, fiber.StatusNotFound, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, tag)
}

// @Summary		Get all
// @Description	Get all tags
// @Tags			tags
// @Accept			json
// @Produce		json
// @Success		200	{array}		domain.Tag
// @Failure		500	{object}	errorResponse
// @Router			/tags [get]
func (h *Handler) tagGetAll(ctx *fiber.Ctx) error {
	tags, err := h.tagService.GetAll(ctx.Context())
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, tags)
}

// @Summary		Get by country ID
// @Description	Get tag by country ID
//
// @Tags			tags
// @Accept			json
// @Produce		json
// @Param			countryID	path		string	true	"Country ID"
// @Success		200			{array}		domain.Tag
// @Failure		400,401,500	{object}	errorResponse
// @Router			/tags/country/{countryID} [get]
func (h *Handler) tagGetAllByCountryID(ctx *fiber.Ctx) error {
	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	tags, err := h.tagService.GetAllByCountryID(ctx.Context(), ctxUser.Settings.CountryID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, tags)
}
