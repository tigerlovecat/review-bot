package config

import "github.com/spf13/viper"

type Redis struct {
	Host     string `json:"host"`
	Password string `json:"password"`
	Db       int    `json:"db"`
}

func InitRedis(cfg *viper.Viper) *Redis {
	return &Redis{
		Host:     cfg.GetString("host"),
		Password: cfg.GetString("password"),
		Db:       cfg.GetInt("db"),
	}
}

var RedisConfig = new(Redis)
