package es

import (
	"context"
	"encoding/json"
	"fmt"
	"log"
	"os"
	"sso/pkg"
	"time"

	"github.com/olivere/elastic/v7"
)

// 查询结果定义
// 查询域名

// 初始化变量
var (
	//	indexName = "elastiflow-*"
	// es集群地址列表：https://sgp.api.es.che/product/es/cluster-list
	//username, password = "anzhihe", "anzhihe"
	client *elastic.Client
	res    *elastic.SearchResult
	err    error
	ctx    context.Context
)

// 流量
type BytesFlow struct {
	Aggregations struct {
		Bytes struct {
			Value float64 `json:"value"`
		} `json:"bytes"`
		DocCount    int    `json:"doc_count"`
		Key         int64  `json:"key"`
		KeyAsString string `json:"key_as_string"`
	} `json:"Aggregations"`
	Key         int64  `json:"Key"`
	KeyAsString string `json:"KeyAsString"`
	KeyNumber   int64  `json:"KeyNumber"`
	DocCount    int    `json:"DocCount"`
}

// 会话
type AgC struct {
	Aggregations struct {
		Num1 struct {
			Value int `json:"value"`
		} `json:"1"`
		Num7 struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Num1 struct {
					Value int `json:"value"`
				} `json:"1"`
				Num4 struct {
					DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
					SumOtherDocCount        int `json:"sum_other_doc_count"`
					Buckets                 []struct {
						Num1 struct {
							Value int `json:"value"`
						} `json:"1"`
						Num8 struct {
							DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
							SumOtherDocCount        int `json:"sum_other_doc_count"`
							Buckets                 []struct {
								Num1 struct {
									Value int `json:"value"`
								} `json:"1"`
								Num6 struct {
									DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
									SumOtherDocCount        int `json:"sum_other_doc_count"`
									Buckets                 []struct {
										Num1 struct {
											Value float64 `json:"value"`
										} `json:"1"`
										Num2 struct {
											Value float64 `json:"value"`
										} `json:"2"`
										Key      string `json:"key"`
										DocCount int    `json:"doc_count"`
									} `json:"buckets"`
								} `json:"6"`
								Key      string `json:"key"`
								DocCount int    `json:"doc_count"`
							} `json:"buckets"`
						} `json:"8"`
						Key      string `json:"key"`
						DocCount int    `json:"doc_count"`
					} `json:"buckets"`
				} `json:"4"`
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"7"`
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"Aggregations"`
	Key         string      `json:"Key"`
	KeyAsString interface{} `json:"KeyAsString"`
	KeyNumber   int         `json:"KeyNumber"`
	DocCount    int         `json:"DocCount"`
}

type FlowData struct {
	Ip      string  `json:"ip"`
	Packets float64 `json:"packets"`
	Bytes   float64 `json:"bytes"`
	Count   int64   `json:"count"`
}
type FlowBytes struct {
	Bytes float64 `json:"bytes"`
	Times int64   `json:"times"`
	Count int64   `json:"count"`
}

type FlowDataC struct {
	Ip       string `json:"ip"`
	Server   string `json:"server"`
	Protocol string `json:"protocol"`
	Packets  int64  `json:"packets"`
	Bytes    int64  `json:"bytes"`
	Count    int64  `json:"count"`
}

type Ag1 struct {
	Aggregations struct {
		Num2 struct {
			Value int `json:"value"`
		} `json:"2"`
		Num5 struct {
			DocCountErrorUpperBound int `json:"doc_count_error_upper_bound"`
			SumOtherDocCount        int `json:"sum_other_doc_count"`
			Buckets                 []struct {
				Num2 struct {
					Value float64 `json:"value"`
				} `json:"2"`
				Num3 struct {
					Value float64 `json:"value"`
				} `json:"3"`
				Key      string `json:"key"`
				DocCount int    `json:"doc_count"`
			} `json:"buckets"`
		} `json:"5"`
		DocCount int    `json:"doc_count"`
		Key      string `json:"key"`
	} `json:"Aggregations"`
	Key         string      `json:"Key"`
	KeyAsString interface{} `json:"KeyAsString"`
	KeyNumber   int         `json:"KeyNumber"`
	DocCount    int         `json:"DocCount"`
}

func ElasticInit() error {
	// 连接es集群
	addr := pkg.Conf().Es.Addrs
	client, err = elastic.NewClient(
		elastic.SetURL(addr...),
		//      elastic.SetBasicAuth(username, password),
		// 允许您指定弹性是否应该定期检查集群（默认为真）
		elastic.SetSniff(false),
		// 设置监控检查时间间隔
		elastic.SetHealthcheckInterval(10*time.Second),
		// 设置错误日志输出
		elastic.SetErrorLog(log.New(os.Stderr, "ELASTIC ", log.LstdFlags)),
		// 设置info日志输出
		elastic.SetInfoLog(log.New(os.Stdout, "", log.LstdFlags)))

	if err != nil {
		fmt.Println("连接失败：%v\n", err)
		return err
	}

	fmt.Println("连接成功", client)
	return nil
}

func QueryGraph(start, end, index int64) ([][]int64, error) {
	//arr := make([]*FlowData, 0)
	arr := make([][]int64, 0)

	ctx = context.Background()

	// g1 := elastic.NewTermsAggregation().Field("network.bytes").OrderByCountDesc().Size(499) //g4

	str1 := time.Unix(int64(start), 0).UTC().Format("2006-01-02T15:04:05Z")
	str2 := time.Unix(int64(end), 0).UTC().Format("2006-01-02T15:04:05Z")

	timeQuery := elastic.NewRangeQuery("@timestamp").
		//              转es的时间
		//      t := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		Gte(str1).
		Lte(str2)
	boolSearch := elastic.NewBoolQuery().Filter(timeQuery)
	fmt.Println(str1)
	fmt.Println(str2)

	aggs := elastic.NewDateHistogramAggregation().
		Field("@timestamp"). // 根据date字段值，对数据进行分组
		//  分组间隔：month代表每月、支持minute（每分钟）、hour（每小时）、day（每天）、week（每周）、year（每年)
		CalendarInterval("minute").
		// 设置返回结果中桶key的时间格式
		Format("yyyy-MM-dd")

	nb := elastic.NewSumAggregation().Field("network.bytes")
	aggs.SubAggregation("bytes", nb)

	indexName := pkg.Conf().Es.IndexName
	res, err = client.Search().Index(indexName).Query(boolSearch).Aggregation("4", aggs).From(0).Pretty(true).Do(ctx) // 执行

	if err != nil {
		return nil, err
	}
	domain, found := res.Aggregations.Terms("4")
	if !found {
		fmt.Println("未找到bucket")
		return nil, err
	}

	for k, bucket := range domain.Buckets {
		b, _ := json.Marshal(bucket)
		log.Printf("====k=%d,v=%s", k, string(b))
		var d BytesFlow
		// var d MyData
		json.Unmarshal(b, &d)
		//              fmt.Println(d)
		temp := d.Aggregations
		t := temp.Key / 1000
		obj := []int64{t, int64(temp.Bytes.Value)}
		arr = append(arr, obj)

		fmt.Printf("times=%d,bytes=%v,count=%d\n", t, temp.Bytes.Value, int(temp.DocCount))
	}

	total := res.TotalHits()
	fmt.Printf("Found %d results\n", total)
	return arr, nil
	//return nil, nil
}
func QueryConversation(start, end, index int64) ([]*FlowDataC, error) {
	fmt.Println("=QueryConversation")
	arr := make([]*FlowDataC, 0)

	ctx = context.Background()

	g1 := elastic.NewTermsAggregation().Field("client.domain").OrderByCountDesc().Size(50)

	g21 := elastic.NewSumAggregation().Field("network.bytes")
	g22 := elastic.NewTermsAggregation().Field("client.ip").OrderByCountDesc().Size(5)

	g31 := elastic.NewSumAggregation().Field("network.bytes")
	g32 := elastic.NewTermsAggregation().Field("server.domain").OrderByCountDesc().Size(25)

	g41 := elastic.NewSumAggregation().Field("network.bytes")
	g42 := elastic.NewTermsAggregation().Field("server.ip").OrderByCountDesc().Size(5)

	g51 := elastic.NewSumAggregation().Field("network.bytes")
	g52 := elastic.NewTermsAggregation().Field("flow.service_name").OrderByCountDesc().Size(10)

	g61 := elastic.NewSumAggregation().Field("network.bytes")
	g62 := elastic.NewSumAggregation().Field("network.packets")

	g1.SubAggregation("1", g21)
	g1.SubAggregation("7", g22)

	g22.SubAggregation("1", g31)
	g22.SubAggregation("4", g32)

	g32.SubAggregation("1", g41)
	g32.SubAggregation("8", g42)

	g42.SubAggregation("1", g51)
	g42.SubAggregation("6", g52)

	g52.SubAggregation("1", g61)
	g52.SubAggregation("2", g62)

	//
	// str1 := "2023-04-03T06:43:00.000Z"
	// str2 := "2023-04-03T07:43:59.999Z"
	str1 := time.Unix(int64(start), 0).UTC().Format("2006-01-02T15:04:05Z")
	str2 := time.Unix(int64(end), 0).UTC().Format("2006-01-02T15:04:05Z")

	timeQuery := elastic.NewRangeQuery("@timestamp").
		Gte(str1).
		Lte(str2)
	boolSearch := elastic.NewBoolQuery().Filter(timeQuery)
	fmt.Println(str1)
	fmt.Println(str2)

	indexName := pkg.Conf().Es.IndexName
	res, err = client.Search().Index(indexName).Query(boolSearch).Aggregation("g1", g1).From(0).Pretty(true).Do(ctx) // 执行

	if err != nil {
		return nil, err
	}

	domain, found := res.Aggregations.Terms("g1")
	if !found {
		fmt.Println("未找到bucket")
		return nil, err
	}
	for _, bucket := range domain.Buckets {
		b, _ := json.Marshal(bucket)
		//fmt.Println(string(b))
		var d AgC
		json.Unmarshal(b, &d)
		//	fmt.Println(d)
		ip := d.Aggregations.Num7.Buckets[0].Key
		for _, v := range d.Aggregations.Num7.Buckets[0].Num4.Buckets {
			server := v.Key
			protocol := v.Num8.Buckets[0].Num6.Buckets[0].Key
			bytes := v.Num8.Buckets[0].Num6.Buckets[0].Num1.Value
			packets := v.Num8.Buckets[0].Num6.Buckets[0].Num2.Value
			obj := &FlowDataC{
				Ip:       ip,
				Server:   server,
				Protocol: protocol,
				Packets:  int64(packets),
				Bytes:    int64(bytes),
				Count:    int64(v.DocCount),
			}
			arr = append(arr, obj)
			//		fmt.Println(ip, v.Key, protocol, bytes, packets)
		}
		// obj := &FlowData{
		//      Ip:      d.Aggregations.Key,
		//      Packets: int64(d.Aggregations.N2.Buckets[0].N4.Value),
		//      Bytes:   int64(d.Aggregations.N1.Value),
		//      Count:   int64(d.Aggregations.DocCount),
		// }
		// arr = append(arr, obj)
		//              fmt.Printf("ip=%s,packet=%f,bytes=%f,count=%d", d.Aggregations.Key, d.Aggregations.N2.Buckets[0].N4.Value, d.Aggregations.N1.Value, d.Aggregations.DocCount)
	}

	total := res.TotalHits()
	fmt.Printf("Found %d results\n", total)
	return arr, nil
}

// 查询日志
func Query1(start, end, index int64) ([]*FlowData, error) {

	arr := make([]*FlowData, 0)

	ctx = context.Background()

	g1 := elastic.NewTermsAggregation().Field("client.domain").OrderByCountDesc().Size(499) //g4

	g22 := elastic.NewSumAggregation().Field("network.bytes")
	g25 := elastic.NewTermsAggregation().Field("client.ip").OrderByCountDesc().Size(5)

	g1.SubAggregation("2", g22)
	g1.SubAggregation("5", g25)

	g32 := elastic.NewSumAggregation().Field("network.bytes")
	g33 := elastic.NewSumAggregation().Field("network.packets")
	g25.SubAggregation("2", g32)
	g25.SubAggregation("3", g33)

	str1 := time.Unix(int64(start), 0).UTC().Format("2006-01-02T15:04:05Z")
	str2 := time.Unix(int64(end), 0).UTC().Format("2006-01-02T15:04:05Z")

	timeQuery := elastic.NewRangeQuery("@timestamp").
		//              转es的时间
		//      t := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		Gte(str1).
		Lte(str2)
	boolSearch := elastic.NewBoolQuery().Filter(timeQuery)
	fmt.Println(str1)
	fmt.Println(str2)

	indexName := pkg.Conf().Es.IndexName
	res, err = client.Search().Index(indexName).Query(boolSearch).Aggregation("4", g1).From(0).Pretty(true).Do(ctx) // 执行

	if err != nil {
		return nil, err
	}
	domain, found := res.Aggregations.Terms("4")
	if !found {
		fmt.Println("未找到bucket")
		return nil, err
	}

	for _, bucket := range domain.Buckets {
		b, _ := json.Marshal(bucket)
		//	fmt.Println(string(b))
		var d Ag1
		// var d MyData
		json.Unmarshal(b, &d)

		temp := d.Aggregations.Num5.Buckets[0]
		obj := &FlowData{
			Ip:      temp.Key,
			Packets: temp.Num3.Value,
			Bytes:   temp.Num2.Value,
			Count:   int64(temp.DocCount),
		}
		arr = append(arr, obj)
	}

	total := res.TotalHits()
	fmt.Printf("Found %d results\n", total)
	return arr, nil
}

func Query2(start, end, index int64) ([]*FlowData, error) {

	arr := make([]*FlowData, 0)

	ctx = context.Background()

	g1 := elastic.NewTermsAggregation().Field("server.domain").OrderByCountDesc().Size(499) //g4

	g22 := elastic.NewSumAggregation().Field("network.bytes")
	g25 := elastic.NewTermsAggregation().Field("server.ip").OrderByCountDesc().Size(5)

	g1.SubAggregation("2", g22)
	g1.SubAggregation("5", g25)

	g32 := elastic.NewSumAggregation().Field("network.bytes")
	g33 := elastic.NewSumAggregation().Field("network.packets")
	g25.SubAggregation("2", g32)
	g25.SubAggregation("3", g33)

	str1 := time.Unix(int64(start), 0).UTC().Format("2006-01-02T15:04:05Z")
	str2 := time.Unix(int64(end), 0).UTC().Format("2006-01-02T15:04:05Z")

	timeQuery := elastic.NewRangeQuery("@timestamp").
		//              转es的时间
		//      t := time.Now().UTC().Format("2006-01-02T15:04:05Z")
		Gte(str1).
		Lte(str2)
	boolSearch := elastic.NewBoolQuery().Filter(timeQuery)
	fmt.Println(str1)
	fmt.Println(str2)

	indexName := pkg.Conf().Es.IndexName
	res, err = client.Search().Index(indexName).Query(boolSearch).Aggregation("4", g1).From(0).Pretty(true).Do(ctx) // 执行

	if err != nil {
		return nil, err
	}
	domain, found := res.Aggregations.Terms("4")
	if !found {
		fmt.Println("未找到bucket")
		return nil, err
	}

	for _, bucket := range domain.Buckets {
		b, _ := json.Marshal(bucket)
		//fmt.Println(string(b))
		var d Ag1
		// var d MyData
		json.Unmarshal(b, &d)
		//              fmt.Println(d)
		temp := d.Aggregations.Num5.Buckets[0]
		obj := &FlowData{
			Ip:      temp.Key,
			Packets: temp.Num3.Value,
			Bytes:   temp.Num2.Value,
			Count:   int64(temp.DocCount),
		}
		arr = append(arr, obj)
		//		fmt.Printf("ip=%s,packet=%f,bytes=%f,count=%d", obj.Ip, obj.Packets, obj.Bytes, obj.Count)
	}

	total := res.TotalHits()
	fmt.Printf("Found %d results\n", total)
	return arr, nil

}
func FormatTimeToInt64(str string) int64 {
	t1, err := time.Parse("2006-01-02T15:04:05Z", str)
	if err != nil {
		fmt.Println("时间解析出错:", err)
		return 0
	}
	return t1.Unix()

}
