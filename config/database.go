package config

import "github.com/spf13/viper"

type Database struct {
	Master string `json:"master"`
	Slave  string `json:"slave"`
}

func InitDatabase(cfg *viper.Viper) *Database {
	return &Database{
		Master: cfg.GetString("master"),
		Slave:  cfg.GetString("slave"),
	}
}

var DatabaseConfig = new(Database)
