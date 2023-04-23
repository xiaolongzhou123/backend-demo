package pkg

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type IConfig struct {
	IP   string `json:"ip"`
	User string `yaml:"User"`
	Pass string `yaml:"Pass"`
	Bs   bool   `yaml:"Bs"`
	Port int    `yaml:"Port"`
}

type Config struct {
	logger    *logrus.Logger
	LogFormat string
	LogLevel  string
	Admin     struct {
		Names []string `yaml:"Names"`
	}
	Pms struct {
		Url    string `yaml:"url"`
		Enable bool   `yaml:"Enable"`
	}
	AES struct {
		KEY string `yaml:"KEY" json:"KEY"`
		IV  string `yaml:"IV" json:"IV"`
	} `yaml:"AES"`
	Es struct {
		IndexName string   `yaml:"indexName"`
		Addrs     []string `yaml:"Addrs"`
	} `yaml:"ES"`
	Nav struct {
		Items []struct {
			Name   string  `yaml:"Name" json:"name"`
			Value  string  `yaml:"Value" json:"value" `
			Id     int     ` json:"id" `
			Config IConfig `yaml:"Config"`
		} `yaml:"Item" json:"item"`
	} `yaml:"Nav"`
	Promethues struct {
		Timeout    string `yaml:"Timeout"`
		AllOut     string `yaml:"AllOut"`
		AllIn      string `yaml:"AllIn"`
		All        string `yaml:"All"`
		QueryRange string `yaml:"QueryRange"`
		Query      string `yaml:"Query"`
		WeekFlow   string `yaml:"WeekFlow"`
		YesToday   struct {
			AllOut string `yaml:"AllOut"`
			AllIn  string `yaml:"AllIn"`
			All    string `yaml:"All"`
		} `yaml:"YesToday"`
		Full struct {
			CPU        string `yaml:"CPU"`
			MEM        string `yaml:"MEM"`
			Temperture string `yaml:"Temperture"`
		} `yaml:"Full"`
	}
	Ldap struct {
		DefaultPass string            `yaml:"DefaultPass"`
		Host        string            `yaml:"Host"`
		BaseDN      string            `yaml:"BaseDN"`
		BindDN      string            `yaml:"BindDN"`
		BinPass     string            `yaml:"BinPass"`
		LdapAttr    []string          `yaml:"LdapAttr"`
		UserDN      map[string]string `yaml:"UserDN"`
	} `yaml:"ldap"`

	Debug     bool   `mapstructure:"Debug" `
	Port      int64  `mapstructure:"Port"`
	JwtExp    int64  `mapstructure:"Jwt_Exp"`
	JwtRef    int64  `mapstructure:"Jwt_Ref"`
	JwtDiff   int64  `mapstructure:"Jwt_Diff"`
	JwtSecret string `mapstructure:"Jwt_Secret"`
}

var config = new(Config)

func Conf() *Config {
	return config
}
func (this *Config) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}
