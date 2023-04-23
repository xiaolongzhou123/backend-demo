package http

import (
	"fmt"
	"sso/cmd/http/admins"
	"sso/cmd/http/middleware"
	"sso/cmd/http/noauth"
	"sso/pkg"

	"github.com/gin-gonic/gin"
)

func RegisterRouter(r *gin.Engine) {
	conf := pkg.Conf()
	//未登陆相关的api

	r.GET("/nav", Nav) //获取ldap的用户
	r.POST("/login", noauth.Login)
	r.GET("/refresh_token", middleware.IsRefresh(), noauth.Refresh_Token)
	r.GET("/ws", Websocket) //web socket不能走认证

	//登陆之后的api，都是认证过的
	r.Use(middleware.IsLogin())
	r.GET("/user", GetUser)            //返回当前用户信息
	r.GET("/ldap_user", GetLdapUser)   //获取ldap的用户
	r.POST("/changepwd", ChangePasswd) //返回当前用户信息

	r.POST("/prome", QueryRange)                                  //查询url请求
	r.POST("/prome/detail", QueryDetail)                          //查询url请求
	r.POST("/prome/bytes/yestoday", QueryYesTodayBytes)           //查询url请求
	r.POST("/prome/full/bytes", QueryFullBytes)                   //查询url请求
	r.POST("/prome/full/todayandyestoday", QueryToDayAndYesToday) //查询url请求

	r.GET("/query/bytes", QueryWeekFlow) //查询url请求
	//websocket调用
	//es
	r.POST("/es/query", QueryES)

	admin := r.Group("/admin")
	//管理员相关api
	if conf.Pms.Enable {
		admin.Use(middleware.IsAdmin())
	}

	admin.POST("/ldap_user", admins.AddLdapUser)   //返回所有用户的数据
	admin.DELETE("/ldap_user", admins.DelLdapUser) //返回所有用户的数据
	admin.PUT("/ldap_user", admins.PutLdapUser)    //返回所有用户的数据

	admin.GET("/users", admins.GetUsers) //返回所有用户的数据

	r.Run(fmt.Sprintf(":%d", conf.Port)) // 监听并在 0.0.0.0:8080 上启动服务
}
