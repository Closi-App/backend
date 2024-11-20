package v1

import (
	"errors"
	"github.com/Closi-App/backend/internal/domain"
	"github.com/Closi-App/backend/internal/utils"
	"github.com/Closi-App/backend/pkg/localizer"
	"github.com/gofiber/fiber/v2"
	"go.mongodb.org/mongo-driver/v2/bson"

	"golang.org/x/text/language"
	"strings"
)

const (
	acceptLanguageHeader = "Accept-Language"
	authorizationHeader  = "Authorization"

	localizerCtxKey = "localizer"
	userCtxKey      = "user"
)

func (h *Handler) getRequestIDFromCtx(ctx *fiber.Ctx) string {
	return ctx.Locals("requestid").(string)
}

func (h *Handler) localizerMiddleware(ctx *fiber.Ctx) error {
	langTag := language.English

	header := ctx.Get(acceptLanguageHeader)
	if header != "" {
		langTags, _, err := language.ParseAcceptLanguage(header)
		if err != nil {
			return h.newResponse(ctx, fiber.StatusInternalServerError, err)
		}

		matcher := language.NewMatcher(h.appSupportedLanguages)
		langTag, _, _ = matcher.Match(langTags...)
	}

	l := h.localizer.SetLanguage(langTag)
	ctx.Locals(localizerCtxKey, l)

	return ctx.Next()
}

func (h *Handler) getLocalizerFromCtx(ctx *fiber.Ctx) *localizer.Localizer {
	l, ok := ctx.Locals(localizerCtxKey).(*localizer.Localizer)
	if !ok {
		l = h.localizer
	}

	return l
}

func (h *Handler) setLocalizerToCtx(ctx *fiber.Ctx, lang string) error {
	langTag, err := utils.ParseLanguage(lang)
	if err != nil {
		return errors.New("error parsing language tag")
	}

	l := h.localizer.SetLanguage(langTag)
	ctx.Locals(localizerCtxKey, l)

	return nil
}

func (h *Handler) userAuthMiddleware(ctx *fiber.Ctx) error {
	header := ctx.Get(authorizationHeader)
	if header == "" {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	headerParts := strings.Split(header, " ")
	if len(headerParts) != 2 || headerParts[0] != "Bearer" {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	accessToken := headerParts[1]

	id, err := h.tokensManager.Parse(accessToken)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	objectID, err := bson.ObjectIDFromHex(id)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	user, err := h.userService.GetByID(ctx.Context(), objectID)
	if err != nil {
		return h.newResponse(ctx, fiber.StatusUnauthorized, domain.ErrUnauthorized)
	}

	ctx.Locals(userCtxKey, user)

	if err = h.setLocalizerToCtx(ctx, user.Settings.Language); err != nil {
		return h.newResponse(ctx, fiber.StatusInternalServerError, err)
	}

	return ctx.Next()
}

func (h *Handler) getUserFromCtx(ctx *fiber.Ctx) (domain.User, error) {
	user, ok := ctx.Locals(userCtxKey).(domain.User)
	if !ok {
		return domain.User{}, errors.New("error getting user from context")
	}

	return user, nil
}
