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
		questions.Get("/", h.questionGetAllWithFilter)
		questions.Get("/:id", h.questionGetByID)

		auth := questions.Group("", h.authUserMiddleware)
		{
			auth.Post("/", h.questionCreate)
			auth.Put("/:id", h.questionUpdate)
			auth.Delete("/:id", h.questionDelete)
		}
	}
}

type questionCreateRequest struct {
	Title          string   `json:"title"`
	Description    string   `json:"description"`
	AttachmentsURL []string `json:"attachments_url"`
	Tags           []string `json:"tags"`
	Points         uint     `json:"points"`
}

// @Summary		Create
// @Description	Create new question
// @Security		UserAuth
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			questionCreateRequest	body		questionCreateRequest	true	"Request"
// @Success		201						{object}	idResponse
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

	if ctxUser.Points < req.Points {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrUserInsufficientPoints)
	}

	id, err := h.questionService.Create(ctx.Context(), service.QuestionCreateInput{
		Title:          req.Title,
		Description:    req.Description,
		AttachmentsURL: req.AttachmentsURL,
		Tags:           req.Tags,
		Points:         req.Points,
		CountryID:      ctxUser.Settings.CountryID,
		UserID:         ctxUser.ID,
	})
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusCreated, idResponse{id.Hex()})
}

// @Summary		Get all with filter
// @Description	Get all question with filter
// @Tags			questions
// @Accept			json
// @Produce		json
// @Param			title		query		string	false	"Question title"
// @Param			tag			query		string	false	"Question tag"
// @Param			countryID	query		string	false	"Country ID"
// @Param			userID		query		string	false	"User ID"
// @Success		200			{array}		domain.Question
// @Failure		400,500		{object}	errorResponse
// @Router			/questions [get]
func (h *Handler) questionGetAllWithFilter(ctx *fiber.Ctx) error {
	title := ctx.Query("title")
	tag := ctx.Query("tag")
	countryID := ctx.Query("country_id")
	userID := ctx.Query("user_id")

	var filter domain.QuestionGetAllFilter

	if title != "" {
		filter.Title = &title
	}
	if tag != "" {
		id, err := bson.ObjectIDFromHex(tag)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}

		filter.Tag = &id
	}
	if countryID != "" {
		id, err := bson.ObjectIDFromHex(countryID)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}

		filter.CountryID = &id
	}
	if userID != "" {
		id, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}

		filter.UserID = &id
	}

	questions, err := h.questionService.GetAll(ctx.Context(), filter)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, questions)
}

// @Summary		Get by ID
// @Description	Get question by ID
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
			return h.newResponse(ctx, fiber.StatusNotFound, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, question)
}

type questionUpdateRequest struct {
	Title          *string  `json:"title"`
	Description    *string  `json:"description"`
	AttachmentsURL []string `json:"attachments_url"`
	Tags           []string `json:"tags"`
	Points         *uint    `json:"points"`
}

// @Summary		Update
// @Description	Update question
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
	var err error

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

	var tags []bson.ObjectID

	if req.Tags != nil {
		for _, tagName := range req.Tags {
			tagID, err := h.tagService.Create(ctx.Context(), service.TagCreateInput{
				Name:      tagName,
				CountryID: ctxUser.Settings.CountryID,
			})
			if err != nil {
				return h.newResponse(ctx, fiber.StatusInternalServerError, err)
			}

			tags = append(tags, tagID)
		}
	}

	if err = h.questionService.Update(ctx.Context(), objectID, ctxUser.ID, domain.QuestionUpdateInput{
		Title:          req.Title,
		Description:    req.Description,
		AttachmentsURL: req.AttachmentsURL,
		Tags:           tags,
		Points:         req.Points,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Delete
// @Description	Delete question
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

	if err = h.questionService.Delete(ctx.Context(), objectID, ctxUser.ID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
