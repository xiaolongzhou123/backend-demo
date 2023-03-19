package http

import (
	"fmt"
	"sso/pkg/jwt"
	"sso/pkg/ldap"
	"sso/pkg/typing"
	"sso/pkg/typing/login"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func GetUser(c *gin.Context) {
	code := 200
	mess := "ok"
	user, _ := c.Get("user")

	rs := typing.NewResp(code, mess, user)
	c.JSON(code, rs)
}

func ChangePasswd(c *gin.Context) {

	var b login.ModifyLogin
	if err := c.ShouldBindWith(&b, binding.JSON); err == nil {

		//连接ldap
		conn, err := ldap.NewLDAP()
		if err != nil {
			c.JSON(200, typing.NewResp(401, "ldap无法连接", err.Error()))
			return
		}
		defer conn.Close()

		user, _ := c.Get("user")
		u, _ := user.(*jwt.Claims)

		if u.User != b.User {
			c.JSON(200, typing.NewResp(10011, "会话用户和提交用户不匹配", struct{}{}))
			return
		}

		//判断用户是否存在,即查询用户dn
		userdn, err := conn.GetUserNameCN(b.User)
		if err != nil {
			c.JSON(200, typing.NewResp(10001, "密码修改失败,获取用户dn失败", err.Error()))
			return
		}

		//判断用户密码是否和原密码相同
		err = conn.LoginV1(userdn, b.OldPass)
		fmt.Println(b.User, b.OldPass)
		if err != nil {
			c.JSON(200, typing.NewResp(10002, "原密码错误,不可修改", struct{}{}))
			return
		}

		//修改密码
		if err = conn.ChangePasswd(userdn, b.OldPass, b.CurPass); err != nil {
			fmt.Println(err)
			c.JSON(200, typing.NewResp(10003, "密码修改失败", err.Error()))
		} else {
			c.JSON(200, typing.NewResp(200, "密码修改成功", struct{}{}))
		}

	} else {
		c.JSON(200, typing.NewResp(401, "post数据不正确", err.Error()))
	}

}
