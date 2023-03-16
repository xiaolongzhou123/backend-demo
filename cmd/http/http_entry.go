package http

import (
	"github.com/gin-gonic/gin"
)

func Start() {
	r := gin.Default()
	ValidInit()       //校验结构体
	RegisterRouter(r) //注册路由

}
