package middleware

import (
	"sso/pkg"
	"sso/pkg/jwt"
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
)

func NotAdmin(mess string, c *gin.Context) {
	code := 200
	rs := typing.NewResp(10001, mess, struct{}{})
	c.JSON(code, rs)
	c.Abort()
}

func IsAdmin() func(c *gin.Context) {
	return func(c *gin.Context) {

		if user, exist := c.Get("user"); exist {
			u, _ := user.(*jwt.Claims)

			names := pkg.Conf().Admin.Names
			ok := false
			for _, name := range names {
				if name == u.User {
					ok = true
					break
				}
			}
			if !ok {
				NotAdmin(u.User+":非管理员", c)
				return
			}
			c.Next()
		}
	}
}
