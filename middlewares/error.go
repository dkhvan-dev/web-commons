package middlewares

import (
	"context"
	"github.com/99designs/gqlgen/graphql"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/constants"
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"github.com/vektah/gqlparser/v2/gqlerror"
	"go.uber.org/zap"
	"net/http"
	"runtime"
	"strings"
)

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		ctx.Next()

		if err, exists := ctx.Get("error"); exists {
			if customErr, ok := err.(*errors.CustomError); ok {
				handleLocalizedError(ctx, customErr)
				return
			}

			handleLocalizedError(ctx, &errors.CustomError{Key: "INTERNAL", Code: http.StatusInternalServerError})
		}
	}
}

func handleLocalizedError(ctx *gin.Context, err *errors.CustomError) {
	lang := ctx.GetHeader(constants.ACCEPT_LANGUAGE)
	if lang == "" {
		lang = "en"
	}

	localizer := i18n.NewLocalizer(config.Bundle, lang)
	message, errLoc := localizer.LocalizeMessage(&i18n.Message{
		ID: err.Key,
	})

	if errLoc != nil {
		message = "An unexpected error occurred"
	}

	logger := config.Logger.With(
		zap.String("error_key", err.Key),
		zap.String("error_msg", message),
		zap.Int("status_code", err.Code),
	)

	if err.File != nil {
		logger.With(zap.String("file", *err.File))
	}

	if err.Line != nil {
		logger.With(zap.Int("line", *err.Line))
	}

	if err.Function != nil {
		logger.With(zap.String("function", *err.Function))
	}

	if err.Key == "INTERNAL" {
		stack := make([]byte, 4096)
		length := runtime.Stack(stack, true)

		logger.With(zap.String("stacktrace", string(stack[:length])))
		logger.Error("Internal error occurred")
	} else {
		logger.Warn("Localized error occurred")
	}

	ctx.JSON(err.Code, gin.H{"error": message})
}

func HandleGraphQLError(ctx context.Context, err error) *gqlerror.Error {
	acceptLang, exists := ctx.Value(constants.ACCEPT_LANGUAGE).(string)
	if !exists {
		acceptLang = "en"
	}

	bugTrack := graphql.DefaultErrorPresenter(ctx, err)
	localizer := i18n.NewLocalizer(config.Bundle, acceptLang)
	whiteSpaceIdx := strings.Index(bugTrack.Message, " ")

	message, errLoc := localizer.LocalizeMessage(&i18n.Message{
		ID: bugTrack.Message[:whiteSpaceIdx],
	})

	if errLoc != nil {
		message = "An unexpected error occurred"
	}

	return gqlerror.Errorf(message)
}
