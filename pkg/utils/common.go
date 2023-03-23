package utils

import (
	"encoding/json"
	"fmt"
	"io/ioutil"
	"net/http"
	"sso/pkg"
)

type Res struct {
	Code    int    `json:"code"`
	Message string `json:"message"`
}

func IsAdmin(username string) bool {
	conf := pkg.Conf()
	url := fmt.Sprintf(conf.Pms.Url, username, "sso_admin")
	fmt.Println("pms", url)
	resp, err := http.Get(url)
	if err != nil {
		// handle error
	}
	defer resp.Body.Close()
	body, err := ioutil.ReadAll(resp.Body)

	var rs Res

	if err := json.Unmarshal(body, &rs); err != nil {
		return false
	}
	if rs.Code == 200 {
		return true
	}

	return false
}
