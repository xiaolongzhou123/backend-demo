package middleware

import (
	"sso/pkg/jwt"
	"sso/pkg/typing"
	"strings"

	"github.com/gin-gonic/gin"
)

func UnLogin(mess string, c *gin.Context) {
	code := 401
	rs := typing.NewResp(code, mess, struct{}{})
	c.JSON(code, rs)
	c.Abort()
}

func IsLogin() func(c *gin.Context) {
	return func(c *gin.Context) {
		// 1. 获取请求头中的 Authorization 认证字段信息
		authStr := c.Request.Header.Get("Authorization")
		// 2.
		if authStr == "" {
			UnLogin("no Authorization", c)
			return
		}
		// 切分, 格式不对，报错
		p := strings.Split(authStr, " ")
		if len(p) != 2 || p[0] != "Bearer" {
			UnLogin("token format err", c)
			return
		}
		// 3. 解析token
		claims, err := jwt.ParseToken(p[1])
		if err != nil {
			UnLogin(err.Error(), c)
			return
		}
		// 保存当前登录的用户信息，存入context 上下文当中去，并继续向下执行
		c.Set("user", claims)
		c.Next()
	}
}
