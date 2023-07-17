package pkg

import (
	"encoding/json"
)

type IConfig struct {
	IP   string `json:"ip,omitempty"`
	User string `json:"user"`
	Pass string `json:"pass"`
	Bs   bool   `json:"bs"`
	Port int    `json:"port"`
}

type Config struct {
	Admins []string `json:"names"`
	Pms    struct {
		Url    string `json:"url"`
		Enable bool   `json:"enable"`
	} `json:"pms"`
	AES struct {
		KEY string `json:"key"`
		IV  string `json:"iv"`
	} `json:"aes"`
	Es struct {
		IndexName string   `json:"indexname"`
		Addrs     []string `json:"addrs"`
	} `json:"es"`
	Nav struct {
		Items []struct {
			Name   string  `json:"name"`
			Value  string  `json:"value"`
			Id     string  `json:"id"`
			Config IConfig `json:"config"`
		} `json:"item"`
	} `json:"nav"`
	Promethues struct {
		Timeout    string `json:"timeout"`
		QueryRange string `json:"queryrange"`
		Query      string `json:"query"`
	} `json:"promethues"`
	Screen struct {
		AllOut   string `json:"allout"`
		AllIn    string `json:"allin"`
		WeekFlow string `json:"weekflow"`
		YesToday struct {
			AllOut string `json:"allout"`
			AllIn  string `json:"allin"`
			All    string `json:"all"`
		} `json:"yestoday"`
		Full struct {
			CPU        string `json:"cpu"`
			MEM        string `json:"mem"`
			Temperture string `json:"temperture"`
		} `json:"full"`
	} `json:"screen"`
	Ldap struct {
		BaseDN   string `json:"basedn"`
		BindPass string `json:"bind_pass"`
		BindDN   string `json:"binddn"`

		DefaultPass string            `json:"defaultpass" `
		Host        string            `json:"host"`
		LdapAttr    []string          `json:"ldapattr" `
		UserDN      map[string]string `json:"userdn"`
	} `json:"ldap"`

	Debug bool  `json:"debug" `
	Port  int64 `json:"port"`
	Jwt   struct {
		Exp    int64  `json:"exp"`
		Ref    int64  `json:"ref"`
		Diff   int64  `json:"diff"`
		Secret string `json:"secret"`
	} `json:"jwt"`
	Ssh struct {
		KeyExchanges []string `json:"keyexchange"`
		Ciphers      []string `json:"ciphers"`
		Macs         []string `json:"macs"`
	}
}

var config = new(Config)

func Conf() *Config {
	return config
}
func (this *Config) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}
