package admins

import (
	"encoding/json"
	"fmt"
	"sso/pkg"
	"sso/pkg/ldap"
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

func PutLdapUser(c *gin.Context) {

	code := 200
	mess := "ok"

	var b typing.LdapUser
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		rs := typing.NewResp(code, err.Error(), struct{}{})
		c.JSON(code, rs)
		return
	}
	fmt.Println(b)

	conn, err := ldap.NewLDAP()
	if err != nil {
		c.JSON(200, typing.NewResp(401, "ldap无法连接", err.Error()))
		return
	}
	defer conn.Close()

	//判断用户是否存在
	userdn, err := conn.GetUserNameCN(b.Cn)
	fmt.Println("==", userdn, err)
	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}
	if userdn == "" {
		c.JSON(200, typing.NewResp(10010, "用户不存在,不能更新", struct{}{}))
		return
	}

	err = conn.ModifyUser(userdn, b)
	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}

	rs := typing.NewResp(code, mess, struct{}{})
	c.JSON(code, rs)
}
func DelLdapUser(c *gin.Context) {
	code := 200
	mess := "ok"

	var b typing.LdapUser
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		rs := typing.NewResp(code, err.Error(), struct{}{})
		c.JSON(code, rs)
		return
	}

	conn, err := ldap.NewLDAP()
	if err != nil {
		c.JSON(200, typing.NewResp(401, "ldap无法连接", err.Error()))
		return
	}
	defer conn.Close()

	//判断用户是否存在
	userdn, err := conn.GetUserNameCN(b.Cn)
	fmt.Println("==", userdn, err)
	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}
	if userdn == "" {
		c.JSON(200, typing.NewResp(10010, "用户不存在,不能更新", struct{}{}))
		return
	}

	err = conn.DelUser(userdn)
	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}

	rs := typing.NewResp(code, mess, struct{}{})
	c.JSON(code, rs)

}
func AddLdapUser(c *gin.Context) {

	conf := pkg.Conf()
	code := 200
	mess := "ok"
	var b typing.LdapUser
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		rs := typing.NewResp(code, err.Error(), struct{}{})
		c.JSON(code, rs)
		return
	}
	b.Sn = b.Cn

	conn, err := ldap.NewLDAP()
	if err != nil {
		c.JSON(200, typing.NewResp(401, "ldap无法连接", err.Error()))
		return
	}
	defer conn.Close()

	//判断用户是否存在
	userdn, err := conn.GetUserNameCN(b.Cn)
	fmt.Println("==", userdn, err)
	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}

	if userdn != "" {
		c.JSON(200, typing.NewResp(10010, "用户已存在", struct{}{}))
		return
	}
	userdn = fmt.Sprintf("cn=%s,%s,%s", b.Cn, conf.Ldap.UserDN["signup"], conf.Ldap.BaseDN)
	fmt.Println("userdn=", userdn)
	arr, _ := conn.GetUsers()
	b.MyId = fmt.Sprintf("%d", len(arr)+1)

	bbb, _ := json.Marshal(b)
	fmt.Println(string(bbb))
	err = conn.AddUser(userdn, b)

	if err != nil {
		c.JSON(200, typing.NewResp(10010, err.Error(), struct{}{}))
		return
	}

	rs := typing.NewResp(code, mess, struct{}{})
	c.JSON(code, rs)

}
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
