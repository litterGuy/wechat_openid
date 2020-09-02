package conf

import (
	"errors"
	"flag"
	"github.com/BurntSushi/toml"
	"log"
)

var (
	confPath string
	Conf     = &Config{}
)

func init() {
	flag.StringVar(&confPath, "confkey", "config.toml", "default config path.")
}

func Init() (err error) {
	if _, err = toml.DecodeFile(confPath, Conf); err != nil {
		log.Println("decode config fail %v", err)
		return
	}
	if err = Conf.Fix(); err != nil {
		return
	}
	return
}

type Config struct {
	AppId     string
	AppSecret string
}

func (c *Config) Fix() (err error) {
	if len(c.AppId) == 0 || len(c.AppSecret) == 0 {
		return errors.New("没有配置appid或者密钥")
	}
	return
}
