package http

import (
	"bytes"
	"fmt"
	"log"
	"net/http"
	"sso/pkg"
	"sso/pkg/aes"
	"sso/pkg/jwt"
	"sso/pkg/myssh"
	"sso/pkg/typing"

	"encoding/base64"
	"encoding/json"

	"github.com/gin-gonic/gin"
	"github.com/gorilla/websocket"
)

var upgrader = websocket.Upgrader{
	ReadBufferSize:  4096,
	WriteBufferSize: 4096,
	// 解决跨域问题
	CheckOrigin: func(r *http.Request) bool {
		//	fmt.Println("CheckOrigin,header===", r.Header)

		return true
	},
}

func Websocket(c *gin.Context) { // 这是基于Gin的Context的实现
	token := c.Query("token")
	fmt.Println("token===", token)

	decoded, err := base64.StdEncoding.DecodeString(token)
	if err != nil {
		code := 401
		rs := typing.NewResp(code, "base64 error or id", struct{}{})
		c.JSON(code, rs)
		return
	}
	decstr := string(decoded)

	txt, err := aes.Decrypt(decstr)
	if err != nil {
		fmt.Println("aes Decrypt error", err)
		code := 401
		rs := typing.NewResp(code, "base64 error or id", struct{}{})
		c.JSON(code, rs)
		return
	}
	fmt.Println("txt======", txt)
	type Data struct {
		Id    int    `json:"id"`
		T     int    `json:"t"`
		Token string `json:"token"`
	}
	var b Data
	if err = json.Unmarshal([]byte(txt), &b); err != nil {
		fmt.Println("websocket unmarshal error:", err, txt)
		code := 401
		rs := typing.NewResp(code, "base64 error or id", struct{}{})
		c.JSON(code, rs)
		return
	}

	claims, err := jwt.ParseToken(b.Token)
	if err != nil {
		fmt.Println("jwt 解析token error", err)
		code := 401
		rs := typing.NewResp(code, "token error", struct{}{})
		c.JSON(code, rs)
		return
	}
	fmt.Println(b.Id, claims)

	// func web(w http.ResponseWriter, r *http.Request) { //这是单独在http下的方法头
	//     conn, _, _, err := ws.UpgradeHTTP(r, w)

	//ip, username, password, Privatekey string, port int
	m := make(map[int]pkg.IConfig, 0)
	items := pkg.Conf().Nav.Items
	for k, v := range items {
		v.Config.IP = v.Value
		m[k] = v.Config
	}
	fmt.Println(m)

	obj, ok := m[b.Id]
	if !ok {
		fmt.Println("id 不在范围")
		code := 401
		rs := typing.NewResp(code, "token error", struct{}{})
		c.JSON(code, rs)
		return
	}

	//client, err := myssh.NewConn("18.18.1.1", "admin", "Nicholas@123", "", true, 51985)

	client, err := myssh.NewConn(obj.IP, obj.User, obj.Pass, "", obj.Bs, obj.Port)

	//client.Enter()
	ws, err := upgrader.Upgrade(c.Writer, c.Request, nil)
	if err != nil {
		log.Println("连接失败：", err)
		return
	}

	defer ws.Close()
	defer client.Close()

	quit := make(chan bool)
	defer close(quit)

	var logBuff = new(bytes.Buffer)
	go client.Recv(ws, logBuff, quit)
	go client.Send(ws, quit)
	client.Wait(quit)
	client.Close()
	fmt.Println("=====end")

}
