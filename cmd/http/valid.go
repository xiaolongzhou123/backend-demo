package http

import (
	"sso/pkg/valid"

	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func ValidInit() {
	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passcheck", valid.Passcheck)
		v.RegisterValidation("usercheck", valid.Usercheck)
	}
}
