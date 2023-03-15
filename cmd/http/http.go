package http

import (
	"fmt"
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

	r.POST("/login", Login)

	r.Use(middleware.IsLogin()) //中间件，这个之后的api。都是认证过的
	r.GET("/user", GetUser)     //返回用户信息

	r.Run(fmt.Sprintf(":%d", conf.Port)) // 监听并在 0.0.0.0:8080 上启动服务
}
