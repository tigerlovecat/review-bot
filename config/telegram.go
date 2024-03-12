package config

import "github.com/spf13/viper"

type Telegram struct {
	Token string `json:"token"`
}

func InitTelegram(cfg *viper.Viper) *Telegram {
	return &Telegram{
		Token: cfg.GetString("token"),
	}
}

var TelegramConfig = new(Telegram)
