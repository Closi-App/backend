package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) initCountryRoutes(router fiber.Router) {
	countries := router.Group("/countries")
	{
		countries.Get("/", h.countryGetAll)
		countries.Get("/:id", h.countryGetByID)

		// TODO: create, delete functions for admins
	}
}

// @Summary		Get all
// @Description	Get all countries
// @Tags			countries
// @Accept			json
// @Produce		json
// @Success		200	{array}		domain.Country
// @Failure		500	{object}	errorResponse
// @Router			/countries [get]
func (h *Handler) countryGetAll(ctx *fiber.Ctx) error {
	countries, err := h.countryService.GetAll(ctx.Context())
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, countries)
}

// @Summary		Get by ID
// @Description	Get country by ID
// @Tags			countries
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Country ID"
// @Success		200			{object}	domain.Country
// @Failure		400,404,500	{object}	errorResponse
// @Router			/countries/{id} [get]
func (h *Handler) countryGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	country, err := h.countryService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrCountryNotFound) {
			return h.newResponse(ctx, fiber.StatusNotFound, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, country)
}
