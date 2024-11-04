package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) initQuestionRoutes(router fiber.Router) {
	questions := router.Group("/questions")
	{
		questions.Get("/", h.questionGet)
		questions.Get("/:id", h.questionGetByID)

		auth := questions.Group("", h.authMiddleware)
		{
			// TODO: getting questions by user id, by location, by tags
			auth.Post("/", h.questionCreate)
			auth.Put("/:id", h.questionUpdate)
			auth.Delete("/:id", h.questionDelete)
		}
	}
}

type questionCreateRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Attachments []string        `json:"attachments"`
	Points      uint            `json:"points"`
	Location    domain.Location `json:"location"`
}

// @Summary		Create question
// @Description	Create a new question
// @Security		UserAuth
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			questionCreateRequest	body		questionCreateRequest	true	"Request"
// @Success		201						{string}	string					"Created"
// @Failure		400,401,500				{object}	errorResponse
// @Router			/questions [post]
func (h *Handler) questionCreate(ctx *fiber.Ctx) error {
	var req questionCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	id, err := h.questionService.Create(ctx.Context(), service.QuestionCreateInput{
		Title:       req.Title,
		Description: req.Description,
		Attachments: req.Attachments,
		Points:      req.Points,
		Location:    req.Location,
		UserID:      ctxUser.ID,
	})
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusCreated, idResponse{id})
}

// @Summary		Get all questions
// @Description	Retrieve a list of all questions
// @Tags			questions
// @Accept			json
// @Produce		json
// @Success		200	{array}		domain.Question
// @Failure		500	{object}	errorResponse
// @Router			/questions [get]
func (h *Handler) questionGet(ctx *fiber.Ctx) error {
	questions, err := h.questionService.GetAll(ctx.Context())
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, questions)
}

// @Summary		Get question by ID
// @Description	Retrieve a specific question by its ID
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Question ID"
// @Success		200			{object}	domain.Question
// @Failure		400,404,500	{object}	errorResponse
// @Router			/questions/{id} [get]
func (h *Handler) questionGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	question, err := h.questionService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrQuestionNotFound) {
			return h.newResponse(ctx, fiber.StatusNotFound, domain.ErrQuestionNotFound)
		}
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, question)
}

type questionUpdateRequest struct {
	Title       string          `json:"title"`
	Description string          `json:"description"`
	Attachments []string        `json:"attachments"`
	Points      uint            `json:"points"`
	Location    domain.Location `json:"location"`
}

// @Summary		Update question
// @Description	Update the details of a specific question by ID
// @Security		UserAuth
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			id						path		string					true	"Question ID"
// @Param			questionUpdateRequest	body		questionUpdateRequest	true	"Request"
// @Success		200						{string}	string					"OK"
// @Failure		400,401,500				{object}	errorResponse
// @Router			/questions/{id} [put]
func (h *Handler) questionUpdate(ctx *fiber.Ctx) error {
	var req questionUpdateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err := h.questionService.Update(ctx.Context(), objectID, ctxUser.ID, service.QuestionUpdateInput{
		Title:       req.Title,
		Description: req.Description,
		Attachments: req.Attachments,
		Points:      req.Points,
		Location:    req.Location,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Delete question
// @Description	Delete a specific question by ID
// @Security		UserAuth
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Question ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/questions/{id} [delete]
func (h *Handler) questionDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err := h.questionService.Delete(ctx.Context(), objectID, ctxUser.ID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
