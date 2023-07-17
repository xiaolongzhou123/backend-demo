package typing

type Promeres struct {
	Status string `json:"status"`
	Query  string `json:"query"`
	Step   string `json:"step"`
	Data   struct {
		ResultType string `json:"resultType"`
		Result     []struct {
			Metric struct {
				IfAlias  string `json:"ifAlias"`
				IfDescr  string `json:"ifDescr"`
				IfIndex  string `json:"ifIndex"`
				IfName   string `json:"ifName"`
				Instance string `json:"instance"`
				SysName  string `json:"sysName"`
				Job      string `json:"job"`
			} `json:"metric"`
			Values [][]interface{} `json:"values"`
		} `json:"result"`
	} `json:"data"`
}

type QueryProme struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Step  int64 `json:"step"`
	Index int64 `json:"index"`
}

type QueryPromeDetail struct {
	Start int64  `json:"start"`
	End   int64  `json:"end"`
	Step  int64  `json:"step"`
	Query string `json:"query"`
}

type QueryBytes struct {
	Start int64 `json:"start"`
	End   int64 `json:"end"`
	Index int   `json:"index"`
	Step  int   `json:"step"`
}

type QueryParam struct {
	Query string `json:"query"`
}
