package admins

import (
	"sso/pkg/ldap"
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
)

func GetUsers(c *gin.Context) {
	code := 200
	mess := "ok"
	var rs *typing.Resp

	conn, err := ldap.NewLDAP()
	if err != nil {
		rs = typing.NewResp(code, err.Error(), struct{}{})
		c.JSON(code, rs)
		return

	}
	defer conn.Close()

	m, err := conn.GetUsers()
	if err != nil {
		rs = typing.NewResp(code, err.Error(), struct{}{})
		c.JSON(code, rs)
		return
	}

	rs = typing.NewResp(code, mess, m)
	c.JSON(code, rs)
}
