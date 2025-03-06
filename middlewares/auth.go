package middlewares

import (
	"errors"
	"github.com/gin-gonic/gin"
	"reflect"
)

func fillAuthor(ctx *gin.Context, request interface{}) error {
	v := reflect.ValueOf(request)
	if v.Kind() != reflect.Ptr || v.Elem().Kind() != reflect.Struct {
		return errors.New("request must be a pointer to a struct")
	}

	v = v.Elem()
	field := v.FieldByName("createdBy")

	if !field.IsValid() {
		return nil
	}

	if !field.CanSet() {
		return errors.New("createdBy field is not updatable")
	}

	fieldValue := reflect.ValueOf(field)
	if field.Type() != fieldValue.Type() {
		return errors.New("provided value does not match with field type")
	}

	field.Set(fieldValue)
	return nil
}
