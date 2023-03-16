package http

import (
	"fmt"
	"sso/cmd/http/admins"
	"sso/cmd/http/middleware"
	"sso/pkg"
	"sso/pkg/valid"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
	"github.com/go-playground/validator/v10"
)

func Start() {
	conf := pkg.Conf()
	r := gin.Default()

	if v, ok := binding.Validator.Engine().(*validator.Validate); ok {
		v.RegisterValidation("passcheck", valid.Passcheck)
		v.RegisterValidation("usercheck", valid.Usercheck)
	}

	//未登陆相关的api
	r.POST("/login", Login)

	//登陆之后的api，都是认证过的
	r.Use(middleware.IsLogin())
	r.GET("/user", GetUser) //返回当前用户信息

	//管理员相关api
	r.Use(middleware.IsAdmin())
	r.GET("/users", admins.GetUsers) //返回所有用户的数据

	r.Run(fmt.Sprintf(":%d", conf.Port)) // 监听并在 0.0.0.0:8080 上启动服务
}
