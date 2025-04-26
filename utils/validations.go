package utils

import (
	"github.com/dkhvan-dev/web-commons/errors"
	"github.com/gin-gonic/gin"
	"github.com/shopspring/decimal"
	"go.mongodb.org/mongo-driver/bson/primitive"
	"reflect"
	"strings"
	"time"
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

func ToString(value interface{}) *string {
	if value == nil {
		return nil
	}

	if str, ok := value.(string); ok {
		return &str
	}

	return nil
}

func ToInt(value interface{}) int {
	if value == nil {
		return 0
	}
	switch v := value.(type) {
	case int:
		return v
	case int32:
		return int(v)
	case int64:
		return int(v)
	case float64:
		return int(v)
	case float32:
		return int(v)
	case primitive.Decimal128:
		i, _, _ := v.BigInt()
		return int(i.Int64())
	default:
		return 0
	}
}

func ToBool(value interface{}) bool {
	if value == nil {
		return false
	}
	if b, ok := value.(bool); ok {
		return b
	}
	return false
}

func ToTime(value interface{}) time.Time {
	if value == nil {
		return time.Time{}
	}
	switch v := value.(type) {
	case primitive.DateTime:
		return v.Time()
	case time.Time:
		return v
	default:
		return time.Time{}
	}
}

func ToBigDecimal(value interface{}) decimal.Decimal {
	if value == nil {
		return decimal.Zero
	}
	switch v := value.(type) {
	case float64:
		return decimal.NewFromFloat(v)
	case float32:
		return decimal.NewFromFloat32(v)
	case int:
		return decimal.NewFromInt(int64(v))
	case int64:
		return decimal.NewFromInt(v)
	case primitive.Decimal128:
		d, _ := decimal.NewFromString(v.String())
		return d
	case string:
		d, err := decimal.NewFromString(v)
		if err != nil {
			return decimal.Zero
		}
		return d
	default:
		return decimal.Zero
	}
}
