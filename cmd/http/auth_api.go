package http

import (
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
)

func GetUser(c *gin.Context) {
	code := 200
	mess := "ok"
	user, _ := c.Get("user")

	rs := typing.NewResp(code, mess, user)
	c.JSON(code, rs)
}
