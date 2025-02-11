package errors

import (
	"fmt"
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

func NewCustomError(key string, code int, r *http.Request) *CustomError {
	pc, file, line, ok := runtime.Caller(2)
	var funcName *string

	if ok {
		fn := runtime.FuncForPC(pc).Name()
		funcName = &fn
	}

	return &CustomError{
		Key:      key,
		Code:     code,
		Request:  r,
		File:     &file,
		Line:     &line,
		Function: funcName,
	}
}

func NotFoundError(key string, r *http.Request) *CustomError {
	return NewCustomError(key, http.StatusNotFound, r)
}

func BadRequestError(key string, r *http.Request) *CustomError {
	return NewCustomError(key, http.StatusBadRequest, r)
}

func ValidationError(key string, r *http.Request) *CustomError {
	return NewCustomError(key, http.StatusBadRequest, r)
}

func AccessDeniedError(r *http.Request) *CustomError {
	return NewCustomError("FORBIDDEN", http.StatusForbidden, r)
}

func UnauthorizedError(r *http.Request) *CustomError {
	return NewCustomError("UNAUTHORIZED", http.StatusUnauthorized, r)
}
