package http

import (
	"sso/pkg"
	"sso/pkg/typing"

	"github.com/gin-gonic/gin"
)

func Nav(c *gin.Context) {
	code := 200
	mess := "ok"

	conf := pkg.Conf()

	type Item struct {
		Name  string `yaml:"Name" json:"name"`
		Value string `yaml:"Value" json:"value" `
		Id    int    ` json:"id" `
	}

	arr := make([]Item, 0)
	for k, v := range conf.Nav.Items {
		item := Item{
			Name:  v.Name,
			Id:    k,
			Value: v.Value,
		}
		arr = append(arr, item)
	}

	c.JSON(code, typing.NewResp(code, mess, arr))
}
