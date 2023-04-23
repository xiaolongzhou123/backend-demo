package http

import (
	"encoding/json"
	"fmt"
	"sso/pkg"
	"sso/pkg/es"
	"sso/pkg/pool"
	"sso/pkg/typing"
	"sso/pkg/utils"
	"time"

	"github.com/gin-gonic/gin"
	"github.com/gin-gonic/gin/binding"
)

//range接口
func QueryWeekFlow(c *gin.Context) {
	code := 200
	mess := "ok"

	conf := pkg.Conf()
	url := conf.Promethues.Query
	query := conf.Promethues.WeekFlow

	rs, err := utils.Query(url, query) //url中，已经有"%s了。"
	if err != nil {
		c.JSON(200, typing.NewResp(10002, err.Error(), struct{}{}))
		return
	}
	c.JSON(code, typing.NewResp(code, mess, rs))
}

func QueryES(c *gin.Context) {

	var b typing.QueryEsInfo
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(200, typing.NewResp(200, err.Error(), struct{}{}))
		return
	}

	bb, _ := json.Marshal(b)
	fmt.Println(string(bb))
	if err := es.ElasticInit(); err != nil {
		c.JSON(200, typing.NewResp(200, err.Error(), struct{}{}))
		return
	}

	if b.Index == 1 {
		rs1, err := es.Query1(b.Start, b.End, b.Index)
		if err != nil {
			c.JSON(200, typing.NewResp(10002, err.Error(), struct{}{}))
			return
		}
		c.JSON(200, typing.NewResp(200, "ok", rs1))
	} else if b.Index == 2 {
		rs1, err := es.Query2(b.Start, b.End, b.Index)
		if err != nil {
			c.JSON(200, typing.NewResp(10002, err.Error(), struct{}{}))
			return
		}
		c.JSON(200, typing.NewResp(200, "ok", rs1))
	} else {

		rs1, err := es.Query1(b.Start, b.End, b.Index)
		if err != nil {
			c.JSON(200, typing.NewResp(10002, err.Error(), struct{}{}))
			return
		}
		c.JSON(200, typing.NewResp(200, "ok", rs1))
	}

}
func QueryToDayAndYesToday(c *gin.Context) {
	code := 200
	mess := "ok"

	var b typing.QueryBytes
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}

	conf := pkg.Conf()
	url := conf.Promethues.QueryRange
	timeout := conf.Promethues.Timeout

	start := time.Unix(b.Start, 0)
	end := time.Unix(b.End, 0)

	diff := (b.End - b.Start) / 60 //转分钟

	Step := utils.FindIndex(diff)

	query := conf.Promethues.Full.CPU

	step := fmt.Sprintf("%d", Step)
	if b.Index == 2 {
		query = conf.Promethues.Full.MEM
	} else if b.Index == 3 {
		query = conf.Promethues.Full.Temperture
	}

	fmt.Println("|||==========", query, start, end, (b.End - b.Start), step)

	tp, _ := time.ParseDuration("-24h")
	start1 := start.Add(tp)
	end1 := end.Add(tp)
	arr := [2][2]time.Time{{start, end}, {start1, end1}}

	//	rs.Query = query
	var obj *typing.Promeres

	p := pool.NewPool(2)
	go func() {
		defer p.Close() //必须调用
		for i := 0; i < 2; i++ {
			ttt := arr[i]
			job := pool.NewJob(i, func() {
				rs, _ := utils.QueryRange(url, query, step, timeout, ttt[0], ttt[1])
				if obj == nil {
					obj = rs
				} else {
					obj.Data.Result = append(obj.Data.Result, rs.Data.Result...)
				}

			})
			p.EntryPoint <- *job
		}
	}()
	p.Run()

	c.JSON(code, typing.NewResp(code, mess, obj))
}
func QueryFullBytes(c *gin.Context) {

	code := 200
	mess := "ok"

	var b typing.QueryBytes
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}

	conf := pkg.Conf()
	url := conf.Promethues.QueryRange
	timeout := conf.Promethues.Timeout

	start := time.Unix(b.Start, 0)
	end := time.Unix(b.End, 0)

	diff := (b.End - b.Start) / 60 //转分钟

	Step := utils.FindIndex(diff)

	query := conf.Promethues.Full.CPU

	step := fmt.Sprintf("%d", Step)
	if b.Index == 2 {
		query = conf.Promethues.Full.MEM
	} else if b.Index == 3 {
		query = conf.Promethues.Full.Temperture
	}

	fmt.Println("|||==========", query, start, end, (b.End - b.Start), step)

	rs, err := utils.QueryRange(url, query, step, timeout, start, end)
	if err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}
	rs.Query = query
	rs.Step = step

	c.JSON(code, typing.NewResp(code, mess, rs))
}
func QueryYesTodayBytes(c *gin.Context) {
	code := 200
	mess := "ok"

	var b typing.QueryBytes
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}

	conf := pkg.Conf()
	url := conf.Promethues.QueryRange
	timeout := conf.Promethues.Timeout

	start := time.Unix(b.Start, 0)
	end := time.Unix(b.End, 0)

	diff := (b.End - b.Start) / 60 //转分钟

	Step := utils.FindIndex(diff)

	query := conf.Promethues.YesToday.All

	step := fmt.Sprintf("%d", Step)
	if b.Index == 1 {
		query = conf.Promethues.YesToday.AllOut
	} else if b.Index == 2 {
		query = conf.Promethues.YesToday.AllIn
	}

	fmt.Println("|||==========", "QueryYesTodayBytes index", b.Index, query, start, end, (b.End - b.Start), step)

	rs, err := utils.QueryRange(url, query, step, timeout, start, end)
	if err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}
	rs.Query = query
	rs.Step = step

	c.JSON(code, typing.NewResp(code, mess, rs))
}
func QueryDetail(c *gin.Context) {
	code := 200
	mess := "ok"

	var b typing.QueryPromeDetail
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}

	conf := pkg.Conf()
	url := conf.Promethues.QueryRange
	timeout := conf.Promethues.Timeout

	start := time.Unix(b.Start, 0)
	end := time.Unix(b.End, 0)

	diff := (b.End - b.Start) / 60 //转分钟
	b.Step = utils.FindIndex(diff)

	step := fmt.Sprintf("%d", b.Step)
	query := b.Query

	fmt.Println("|||==========", query, start, end, (b.End - b.Start), step)

	rs, err := utils.QueryRange(url, query, step, timeout, start, end)
	if err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}
	rs.Query = query
	rs.Step = step

	c.JSON(code, typing.NewResp(code, mess, rs))
}
func QueryRange(c *gin.Context) {
	code := 200
	mess := "ok"

	var b typing.QueryProme
	if err := c.ShouldBindWith(&b, binding.JSON); err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}

	conf := pkg.Conf()
	url := conf.Promethues.QueryRange
	timeout := conf.Promethues.Timeout

	start := time.Unix(b.Start, 0)
	end := time.Unix(b.End, 0)
	step := fmt.Sprintf("%d", b.Step)
	query := ""
	if b.Index == 1 {
		query = conf.Promethues.AllOut
	} else {
		query = conf.Promethues.AllIn
	}

	fmt.Println("====", query, start, end, step)

	rs, err := utils.QueryRange(url, query, step, timeout, start, end)
	if err != nil {
		c.JSON(code, typing.NewResp(code, err.Error(), struct{}{}))
		return
	}
	rs.Query = query
	rs.Step = step

	c.JSON(code, typing.NewResp(code, mess, rs))
}
