package errors

import (
	"fmt"
	"github.com/dkhvan-dev/web-commons/config"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
	"net/http"
	"runtime"
)

type CustomError struct {
	Key      string
	Code     int
	Request  *http.Request
	File     *string
	Line     *int
	Function *string
}

func (e *CustomError) Error() string {
	location := ""

	if e.File != nil && e.Line != nil && e.Function != nil {
		location = fmt.Sprintf(" (file: %s, line: %d, function: %s)", *e.File, *e.Line, *e.Function)
	}

	return fmt.Sprintf("%s%s", e.Key, location)
}

func NewCustomError(key string, code int, ctx *gin.Context) *CustomError {
	pc, file, line, ok := runtime.Caller(2)
	var funcName *string

	if ok {
		fn := runtime.FuncForPC(pc).Name()
		funcName = &fn
	}

	err := &CustomError{
		Key:      key,
		Code:     code,
		File:     &file,
		Line:     &line,
		Function: funcName,
	}

	if ctx != nil {
		err.Request = ctx.Request
		ctx.Set("error", err)
	}

	return err
}

func NotFoundError(key string, ctx *gin.Context) *CustomError {
	return NewCustomError(key, http.StatusNotFound, ctx)
}

func BadRequestError(key string, ctx *gin.Context) *CustomError {
	return NewCustomError(key, http.StatusBadRequest, ctx)
}

func ValidationError(key string, ctx *gin.Context) *CustomError {
	return NewCustomError(key, http.StatusBadRequest, ctx)
}

func AccessDeniedError(ctx *gin.Context) *CustomError {
	return NewCustomError("FORBIDDEN", http.StatusForbidden, ctx)
}

func UnauthorizedError(ctx *gin.Context) *CustomError {
	return NewCustomError("UNAUTHORIZED", http.StatusUnauthorized, ctx)
}

func HandleInternalError(err error) {
	if err == nil {
		return
	}

	stack := make([]byte, 4096)
	length := runtime.Stack(stack, true)

	config.Logger.Error("Internal error occurred",
		zap.Error(err),
		zap.String("stacktrace", string(stack[:length])),
	)
}
