package noauth

import (
	"fmt"
	"sso/pkg/jwt"
	"strconv"

	"sso/pkg/ldap"
	"sso/pkg/typing"
	"sso/pkg/typing/login"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	//"github.com/go-ldap/ldap"
)

func Refresh_Token(c *gin.Context) {
	useri, _ := c.Get("user")

	claims, _ := useri.(*jwt.Claims)
	token, refresh := jwt.CreateToken(claims.Uid, claims.User, claims.CN, claims.Mail)

	rs := make(map[string]string, 0)
	rs["access_token"] = fmt.Sprintf("%s %s", "Bearer", token)
	rs["refresh_token"] = fmt.Sprintf("%s %s", "Bearer", refresh)
	c.JSON(200, typing.NewResp(200, "刷新token成功", rs))
}
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

		cn, _ := m["cn"].(string)
		mail, _ := m["mail"].(string)
		uidstr, _ := m["uid"].(string)
		uid, err := strconv.ParseInt(uidstr, 10, 64)

		//如果存在，判断密码是否正确
		err, reason := conn.Login(b.User, b.Pass)
		if err != nil {
			c.JSON(200, typing.NewResp(401, "登陆失败", reason))
			return
		}

		token, refresh := jwt.CreateToken(uid, b.User, cn, mail)

		rs := make(map[string]string, 0)
		rs["access_token"] = fmt.Sprintf("%s %s", "Bearer", token)
		rs["refresh_token"] = fmt.Sprintf("%s %s", "Bearer", refresh)

		c.JSON(200, typing.NewResp(200, "登陆成功", rs))
	} else {
		c.JSON(200, typing.NewResp(401, "post数据不正确", struct{}{}))
	}
}
