package utils

import (
	"encoding/base32"
	"encoding/json"
	"fmt"
	"io"
	"io/ioutil"
	"log"
	"math/rand"
	"net/http"
	"net/url"
	"sso/pkg"
	"sso/pkg/typing"
	"strconv"
	"strings"
	"time"
)

type Res struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func IsAdmin(username, auth string) bool {
	conf := pkg.Conf()
	url := fmt.Sprintf(conf.Pms.Url, username, "sso_admin")

	client := &http.Client{
		Timeout: 5 * time.Second,
	}

	req, err := http.NewRequest("GET", url, nil)
	if err != nil {
		log.Println(err)
		return false
	}
	req.Header.Set("Authorization", auth)

	resp, err := client.Do(req)
	if err != nil {
		log.Println(err)
		return false
	}
	defer resp.Body.Close()

	body, err := io.ReadAll(resp.Body)
	if err != nil {
		log.Println(err)
		return false
	}

	var rs Res
	if err := json.Unmarshal(body, &rs); err != nil {
		return false
	}
	if rs.Code == 200 {
		return true
	}

	return false
}
func GetRandomString(length int) string {
	// length := 10
	str := "0123456789abcdefghijklmnopqrstuvwxyzABCDEFGHIJKLMNOPQRSTUVWXYZ"
	bytes := []byte(str)
	result := []byte{}
	r := rand.New(rand.NewSource(time.Now().UnixNano()))
	for i := 0; i < length; i++ {
		result = append(result, bytes[r.Intn(len(bytes))])
	}
	return string(result)

}
func CreatePrivKey() string {
	data := GetRandomString(10)
	str := base32.StdEncoding.EncodeToString([]byte(data))
	return str
}

func formatTime(t time.Time) string {
	return strconv.FormatFloat(float64(t.Unix())+float64(t.Nanosecond())/1e9, 'f', -1, 64)
}

// cur := time.Now()
// start := cur.Add(-time.Hour)
// end := cur
// query := "irate(ifOutOctets{job='AR',ifName='GigabitEthernet0/0/6'}[2m])"
// addr := "http://18.18.2.2:9090/api/v1/query_range"
// rs, err := QueryRange(addr, query, "60", "5s", start, end)
func QueryRange(addr, query, step, timeout string, start, end time.Time) (*typing.Promeres, error) {
	u := url.URL{}
	q := u.Query()
	q.Set("query", query)
	q.Set("start", formatTime(start))
	q.Set("end", formatTime(end))
	q.Set("step", step)
	q.Set("timeout", timeout)
	rs := &typing.Promeres{}

	req, err := http.NewRequest("POST", addr, strings.NewReader(q.Encode()))
	if err != nil {
		fmt.Println("err:=", err)
		return rs, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["Idempotency-Key"] = nil

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("do err:", err)
		return rs, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return rs, err
	}
	fmt.Println(string(body))
	json.Unmarshal(body, rs)
	return rs, nil
}

func Query(addr, query string) (*typing.Promeres, error) {
	u := url.URL{}
	q := u.Query()
	q.Set("query", query)
	rs := &typing.Promeres{}
	nurl := fmt.Sprintf("%s?%s", addr, q.Encode())

	req, err := http.NewRequest("GET", nurl, nil)
	if err != nil {
		fmt.Println("err:=", err)
		return rs, err
	}
	req.Header.Set("Content-Type", "application/x-www-form-urlencoded")
	req.Header["Idempotency-Key"] = nil

	client := &http.Client{}
	response, err := client.Do(req)
	if err != nil {
		fmt.Println("do err:", err)
		return rs, err
	}
	defer response.Body.Close()

	body, err := ioutil.ReadAll(response.Body)
	if err != nil {
		fmt.Println("ioutil.ReadAll err:", err)
		return rs, err
	}
	fmt.Println(string(body))
	json.Unmarshal(body, rs)
	return rs, nil
}
func FindIndex(diff int64) int64 {
	if diff > 8*7*24*60 { //8周，每周7天 每天24小时 每小时60分钟
		return 19353
	} else if diff > 4*7*24*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 9676
	} else if diff > 2*7*24*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 4838
	} else if diff > 1*7*24*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 2419
	} else if diff > 2*24*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 691
	} else if diff > 1*24*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 345
	} else if diff > 12*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 172
	} else if diff > 6*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 86
	} else if diff > 2*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 28
	} else if diff > 1*60 { //4周，每周7天 每天24小时 每小时60分钟
		return 14
	} else if diff > 7 { //4周，每周7天 每天24小时 每小时60分钟
		return 30
	} else if diff > 3 { //4周，每周7天 每天24小时 每小时60分钟
		return 15
	}
	return 1
}
