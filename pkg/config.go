package pkg

import (
	"encoding/json"

	"github.com/sirupsen/logrus"
)

type Config struct {
	logger    *logrus.Logger
	LogFormat string
	LogLevel  string

	Debug bool  `mapstructure:"DEBUG" yaml:"-"`
	Port  int64 `mapstructure:"PORT"`
}

var config = new(Config)

func Conf() *Config {
	return config
}
func (this *Config) String() string {
	b, _ := json.Marshal(this)
	return string(b)
}
