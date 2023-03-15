package middleware

import (
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
)

func UnLogin(c *gin.Context) {
	code := 401
	rs := typing.NewResp(code, "noauth", struct{}{})
	c.JSON(code, rs)
}
func IsLogin() gin.HandlerFunc {
	return func(c *gin.Context) {
	}
}
