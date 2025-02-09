package errors

import (
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/dkhvan-dev/web-commons/constants"
	"github.com/gin-gonic/gin"
	"github.com/nicksnyder/go-i18n/v2/i18n"
	"net/http"
)

type CustomError struct {
	Key     string
	Code    int
	Request *http.Request
}

func (e *CustomError) Error() string {
	return e.Key
}

func NotFoundError(key string, r *http.Request) *CustomError {
	return &CustomError{Key: key, Code: 404, Request: r}
}

func BadRequestError(key string, r *http.Request) *CustomError {
	return &CustomError{Key: key, Code: 400, Request: r}
}

func ValidationError(key string, r *http.Request) *CustomError {
	return &CustomError{Key: key, Code: 400, Request: r}
}

func AccessDeniedError(r *http.Request) *CustomError {
	return &CustomError{Key: "FORBIDDEN", Code: 403, Request: r}
}

func UnauthorizedError(r *http.Request) *CustomError {
	return &CustomError{Key: "UNAUTHORIZED", Code: 401, Request: r}
}

func handleLocalizedError(ctx *gin.Context, err *CustomError) {
	lang := ctx.GetHeader(constants.ACCEPT_LANGUAGE)
	if lang == "" {
		lang = "en"
	}

	localizer := i18n.NewLocalizer(config.Bundle, lang)
	message := localizer.MustLocalize(&i18n.LocalizeConfig{
		MessageID: err.Key,
	})

	ctx.JSON(err.Code, gin.H{"error": message})
}

func ErrorHandler() gin.HandlerFunc {
	return func(ctx *gin.Context) {
		defer func() {
			if rec := recover(); rec != nil {
				if err, ok := rec.(*CustomError); ok {
					handleLocalizedError(ctx, err)
					return
				}

				ctx.JSON(http.StatusInternalServerError, gin.H{"error": "internal server error"})
			}
		}()
		ctx.Next()
	}
}
