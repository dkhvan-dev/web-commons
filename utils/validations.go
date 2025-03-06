package utils

import (
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"reflect"
	"strings"
)

func Validate(s interface{}, ctx *gin.Context) *errors.CustomError {
	v := reflect.ValueOf(s)
	if v.Kind() == reflect.Ptr {
		v = v.Elem()
	}

	t := v.Type()
	for i := 0; i < v.NumField(); i++ {
		field := v.Field(i)
		fieldType := t.Field(i)
		validateTag := fieldType.Tag.Get("validate")

		if strings.Contains(validateTag, "isEmpty") {
			if IsEmpty(field) {
				return errors.ValidationError("FIELD_IS_EMPTY", ctx)
			}
		}
	}

	return nil
}

func IsEmpty(v reflect.Value) bool {
	switch v.Kind() {
	case reflect.Ptr, reflect.Interface:
		return v.IsNil()
	case reflect.String:
		return v.String() == ""
	case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
		return v.Int() == 0
	case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
		return v.Uint() == 0
	case reflect.Float32, reflect.Float64:
		return v.Float() == 0
	default:
		return !v.IsValid()
	}
}
