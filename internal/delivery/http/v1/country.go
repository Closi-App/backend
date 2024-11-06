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
		countries.Get("/", h.countryGet)
		countries.Get("/:id", h.countryGetByID)

		//auth := countries.Group("") TODO: admins middleware
		//{
		//	auth.Post("/", h.countryCreate)
		//}
	}
}

// @Summary		Get all
// @Description	Get all countries
// @Tags		countries
// @Accept		json
// @Produce		json
// @Success		200	{array}	domain.Country
// @Failure		500	{object} errorResponse
// @Router		/countries [get]
func (h *Handler) countryGet(ctx *fiber.Ctx) error {
	countries, err := h.countryService.Get(ctx.Context())
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, countries)
}

// @Summary		Get by ID
// @Description	Get country by ID
// @Tags		countries
// @Accept		json
// @Produce		json
// @Param		id path	string true	"Country ID"
// @Success		200	{object} domain.Country
// @Failure		400,404,500	{object} errorResponse
// @Router		/countries/{id} [get]
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

//type countryCreateRequest struct {
//	Name map[domain.Language]string `json:"name"`
//}
//
//func (h *Handler) countryCreate(ctx *fiber.Ctx) error {
//	var req countryCreateRequest
//	if err := ctx.BodyParser(&req); err != nil {
//		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
//	}
//
//	id, err := h.countryService.Create(ctx.Context(), service.CountryCreateInput{
//		Name: req.Name,
//	})
//	if err != nil {
//		return h.newResponse(ctx, fiber.StatusInternalServerError, domain.ErrInternalServerError)
//	}
//
//	return h.newResponse(ctx, fiber.StatusCreated, idResponse{id.Hex()})
//}
