package valid

import (
	"fmt"

	"github.com/go-playground/validator/v10"
)

var Passcheck validator.Func = func(fl validator.FieldLevel) bool {
	fmt.Println("pass+++++++++++++++++++++++++++++++++++++}")
	data, ok := fl.Field().Interface().(string)
	fmt.Println("=================================Passcheck validator.Func:", data, ok)
	return true
}
var Usercheck validator.Func = func(fl validator.FieldLevel) bool {
	fmt.Println("user+++++++++++++++++++++++++++++++++++++}")
	data, ok := fl.Field().Interface().(string)
	fmt.Println("=================================Passcheck validator.Func:", data, ok)
	return true
}

func Test() {
}
