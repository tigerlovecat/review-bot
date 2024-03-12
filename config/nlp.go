package config

import "github.com/spf13/viper"

type Nlp struct {
	Key    string `json:"key"`
	Secret string `json:"secret"`
}

func InitNlp(cfg *viper.Viper) *Nlp {
	return &Nlp{
		Key:    cfg.GetString("key"),
		Secret: cfg.GetString("secret"),
	}
}

var NlpConfig = new(Nlp)
