package config

import (
	"github.com/spf13/viper"
)

type Application struct {
	ReadTimeout   int
	WriterTimeout int
	Mode          string
	Name          string
	Host          string
	Port          string
	Prefix        string
}

func InitApplication(cfg *viper.Viper) *Application {
	return &Application{
		ReadTimeout:   cfg.GetInt("read_timeout"),
		WriterTimeout: cfg.GetInt("write_timeout"),
		Host:          cfg.GetString("host"),
		Port:          cfg.GetString("port"),
		Name:          cfg.GetString("name"),
		Mode:          cfg.GetString("mode"),
		Prefix:        cfg.GetString("prefix"),
	}
}

var ApplicationConfig = new(Application)
