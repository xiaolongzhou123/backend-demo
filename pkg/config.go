package pkg

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Config struct {
	logger    *logrus.Logger
	LogFormat string
	LogLevel  string
	Admin     struct {
		Names []string `yaml:"Names"`
	}
	Ldap struct {
		Host     string            `yaml:"Host"`
		BaseDN   string            `yaml:"BaseDN"`
		BindDN   string            `yaml:"BindDN"`
		BinPass  string            `yaml:"BinPass"`
		LdapAttr []string          `yaml:"LdapAttr"`
		UserDN   map[string]string `yaml:"UserDN"`
	} `yaml:"ldap"`

	Debug     bool   `mapstructure:"Debug" `
	Port      int64  `mapstructure:"Port"`
	JwtExp    int64  `mapstructure:"Jwt_Exp"`
	JwtRef    int64  `mapstructure:"Jwt_Ref"`
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
