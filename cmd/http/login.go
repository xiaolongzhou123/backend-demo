package http

import (
	"sso/pkg/typing"
	"sso/pkg/typing/login"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func Login(c *gin.Context) {
	var b login.Login
	if err := c.ShouldBindWith(&b, binding.JSON); err == nil {
		c.JSON(200, typing.NewResp(200, "密码正确", struct{}{}))
	} else {
		c.JSON(200, typing.NewResp(200, "密码正确", struct{}{}))
	}
}
