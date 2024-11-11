package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/service"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"
)

func (h *Handler) initAnswerRoutes(router fiber.Router) {
	answers := router.Group("/answers")
	{
		answers.Get("/", h.answerGetAllWithFilter)
		answers.Get("/:id", h.answerGetByID)

		auth := answers.Group("", h.authUserMiddleware)
		{
			auth.Post("/", h.answerCreate)
			auth.Put("/:id", h.answerUpdate)
			auth.Delete("/:id", h.answerDelete)

			auth.Put("/:id/likes", h.answerAddLike)
			auth.Delete("/:id/likes", h.answerRemoveLike)
		}
	}
}

type answerCreateRequest struct {
	Text       string `json:"text"`
	QuestionID string `json:"question_id"`
}

// @Summary		Create
// @Description	Create new answer
// @Security		UserAuth
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			answerCreateRequest	body		answerCreateRequest	true	"Request"
// @Success		201					{object}	idResponse
// @Failure		400,401,500			{object}	errorResponse
// @Router			/answers [post]
func (h *Handler) answerCreate(ctx *fiber.Ctx) error {
	var req answerCreateRequest
	if err := ctx.BodyParser(&req); err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	questionObjectID, err := bson.ObjectIDFromHex(req.QuestionID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	id, err := h.answerService.Create(ctx.Context(), service.AnswerCreateInput{
		Text:       req.Text,
		QuestionID: questionObjectID,
		UserID:     ctxUser.ID,
	})
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusCreated, idResponse{id.Hex()})
}

// @Summary		Get all with filter
// @Description	Get all answers with filter
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			questionID	query		string	false	"Question ID"
// @Param			userID		query		string	false	"User ID"
// @Success		200			{array}		domain.Question
// @Failure		400,500		{object}	errorResponse
// @Router			/answers [get]
func (h *Handler) answerGetAllWithFilter(ctx *fiber.Ctx) error {
	questionID := ctx.Query("question_id")
	userID := ctx.Query("user_id")

	var filter domain.AnswerGetAllFilter

	if questionID != "" {
		id, err := bson.ObjectIDFromHex(questionID)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}

		filter.QuestionID = &id
	}
	if userID != "" {
		id, err := bson.ObjectIDFromHex(userID)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
		}

		filter.UserID = &id
	}

	answers, err := h.answerService.GetAll(ctx.Context(), filter)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, answers)
}

// @Summary		Get by ID
// @Description	Get answer by ID
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Answer ID"
// @Success		200			{object}	domain.Answer
// @Failure		400,404,500	{object}	errorResponse
// @Router			/answers/{id} [get]
func (h *Handler) answerGetByID(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	answer, err := h.answerService.GetByID(ctx.Context(), objectID)
	if err != nil {
		if errors.Is(err, domain.ErrAnswerNotFound) {
			return h.newResponse(ctx, fiber.StatusNotFound, err)
		}

		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, answer)
}

type answerUpdateRequest struct {
	Text *string `json:"text"`
}

// @Summary		Update
// @Description	Update answer
// @Security		UserAuth
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			id					path		string				true	"Answer ID"
// @Param			answerUpdateRequest	body		answerUpdateRequest	true	"Request"
// @Success		200					{string}	string				"OK"
// @Failure		400,401,500			{object}	errorResponse
// @Router			/answers/{id} [put]
func (h *Handler) answerUpdate(ctx *fiber.Ctx) error {
	var err error

	var req answerUpdateRequest
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

	if err = h.answerService.Update(ctx.Context(), objectID, ctxUser.ID, domain.AnswerUpdateInput{
		Text: req.Text,
	}); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Delete
// @Description	Delete answer
// @Security		UserAuth
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Answer ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/answers/{id} [delete]
func (h *Handler) answerDelete(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	ctxUser, err := h.getUserFromCtx(ctx)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	if err = h.answerService.Delete(ctx.Context(), objectID, ctxUser.ID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Add like
// @Description	Add like for answer
// @Security		UserAuth
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Answer ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/answers/{id}/likes [put]
func (h *Handler) answerAddLike(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	if err = h.answerService.AddLike(ctx.Context(), objectID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}

// @Summary		Remove like
// @Description	Remove like for answer
// @Security		UserAuth
// @Tags			answers
// @Accept			json
// @Produce		json
// @Param			id			path		string	true	"Answer ID"
// @Success		200			{string}	string	"OK"
// @Failure		400,401,500	{object}	errorResponse
// @Router			/answers/{id}/likes [delete]
func (h *Handler) answerRemoveLike(ctx *fiber.Ctx) error {
	id := ctx.Params("id")
	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusBadRequest, domain.ErrBadRequest)
	}

	if err = h.answerService.RemoveLike(ctx.Context(), objectID); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return h.newResponse(ctx, fiber.StatusOK, nil)
}
