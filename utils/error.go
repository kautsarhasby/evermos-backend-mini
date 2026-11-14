package utils

import (
	"reflect"

	"github.com/go-playground/validator/v10"
)

func GetCustomError(err error, obj interface{}) string {
	if ve, ok := err.(validator.ValidationErrors); ok {
		out := reflect.TypeOf(obj).Elem()
		for _, fe := range ve {
			field, _ := out.FieldByName(fe.Field())
			msg := field.Tag.Get("msg")
			if msg != "" {
				return msg
			}
			return fe.Error()
		}
	}
	return err.Error()
}
