package valid

import (
	"github.com/go-playground/validator/v10"
)

var Passcheck validator.Func = func(fl validator.FieldLevel) bool {
	if data, ok := fl.Field().Interface().(string); ok {
		length := len(data)
		if length > 20 || length < 2 {
			return false
		}
		return true
	}
	return false
}
var Usercheck validator.Func = func(fl validator.FieldLevel) bool {
	if data, ok := fl.Field().Interface().(string); ok {

		length := len(data)
		if length < 2 {
			return false
		}

		return true
	}
	return false
}

func Test() {
}
