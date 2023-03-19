package noauth

import (
	"fmt"
	"sso/pkg/jwt"

	"sso/pkg/ldap"
	"sso/pkg/typing"
	"sso/pkg/typing/login"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	//"github.com/go-ldap/ldap"
)

func Login(c *gin.Context) {
	var b login.Login
	if err := c.ShouldBindWith(&b, binding.JSON); err == nil {

		//连接ldap
		conn, _ := ldap.NewLDAP()
		if err != nil {
			c.JSON(200, typing.NewResp(401, "ldap无法连接", err.Error()))
			return
		}
		defer conn.Close()

		//判断用户是否存在
		m, err := conn.GetUser(b.User)
		if err != nil {
			c.JSON(200, typing.NewResp(401, "用户不存在", err.Error()))
			return
		}

		sn, _ := m["sn"].(string)
		mail, _ := m["mail"].(string)

		//如果存在，判断密码是否正确
		err, reason := conn.Login(b.User, b.Pass)
		if err != nil {
			c.JSON(200, typing.NewResp(401, "登陆失败", reason))
			return
		}

		token, _ := jwt.CreateToken(b.User, sn, mail)
		c.JSON(200, typing.NewResp(200, "登陆成功", fmt.Sprintf("%s %s", "Bearer", token)))
	} else {
		c.JSON(200, typing.NewResp(401, "post数据不正确", struct{}{}))
	}
}
