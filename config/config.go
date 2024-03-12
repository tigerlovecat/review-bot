package config

import (
	"github.com/spf13/viper"
	"os"
	"path"
)

var cfgApplication *viper.Viper
var cfgDatabase *viper.Viper
var cfgRedis *viper.Viper
var cfgNlp *viper.Viper
var cfgTelegram *viper.Viper
var ProjectHomePath, _ = os.Getwd()

func Setup(pathDir string) {
	if pathDir == "" {
		pathDir = path.Join(ProjectHomePath, "config/settings.yaml")
	}
	settingCfg := viper.New()
	settingCfg.SetConfigFile(pathDir)
	if err := settingCfg.ReadInConfig(); err != nil {
		panic("找不到配置文件{settings.yml}." + err.Error())
	}
	// 启动参数
	cfgApplication = settingCfg.Sub("application")
	if cfgApplication == nil {
		panic("config not found  application")
	}
	ApplicationConfig = InitApplication(cfgApplication)

	// mysql配置文件初始化
	cfgDatabase = settingCfg.Sub("database.mysql")
	if cfgDatabase == nil {
		panic("config not found database.mysql")
	}
	DatabaseConfig = InitDatabase(cfgDatabase)

	// redis配置文件初始化
	//cfgRedis = settingCfg.Sub("redis")
	//if cfgRedis == nil {
	//	panic("config not fount redis")
	//}
	//RedisConfig = InitRedis(cfgRedis)

	// nlp 配置文件初始化
	cfgNlp = settingCfg.Sub("nlp")
	if cfgNlp == nil {
		panic("config not fount nlp")
	}
	NlpConfig = InitNlp(cfgNlp)

	// telegram 配置文件初始化
	cfgTelegram = settingCfg.Sub("telegram")
	if cfgTelegram == nil {
		panic("config not fount telegram")
	}
	TelegramConfig = InitTelegram(cfgTelegram)

}
