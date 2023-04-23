package middleware

import (
	"sso/pkg/jwt"
	"sso/pkg/typing"
	"sso/pkg/utils"

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

			iauth, _ := c.Get("Authorization")
			auth, _ := iauth.(string)
			ok := utils.IsAdmin(u.CN, auth)
			if !ok {
				NotAdmin(u.CN+":非管理员", c)
				return
			}

			c.Next()
		}
	}
}
