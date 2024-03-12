package cmd

import (
	"context"
	"errors"
	"fmt"
	"log"
	"net/http"
	"os"
	"os/signal"
	"review-bot/database"
	"time"

	conf "review-bot/config"
	"review-bot/router"

	"github.com/spf13/cobra"
)

var (
	config   string
	port     string
	mode     string
	StartCmd = &cobra.Command{
		Use:     "server",
		Short:   "启动服务",
		Example: "./wb server config/settings.yaml",
		PreRun: func(cmd *cobra.Command, args []string) {
			usage()
			setup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return runServer()
		},
	}
)

func init() {
	StartCmd.PersistentFlags().StringVarP(&config, "config", "c", "config/settings.yaml", "Start server with provided configuration file")
	StartCmd.PersistentFlags().StringVarP(&port, "port", "p", "5000", "Tcp port server listening on")
	StartCmd.PersistentFlags().StringVarP(&mode, "mode", "m", "dev", "server mode ; eg:dev,test,prod")
}

func usage() {
	usageStr := `starting api server`
	log.Printf("%s\n", usageStr)
}

func setup() {
	// 初始化配置文件
	conf.Setup(config)
	// 初始化mysql
	database.MysqlSetup()
	// init redis
	//database.RedisSetup()
	// sop cron
	//if conf.ApplicationConfig.Mode == "debug" {
	//	task.StartTask()
	//}
	// 执行sop定时任务
	//task.Init()
}

func runServer() error {
	fmt.Println("服务器正在启动中...123456")
	r := router.InitRouter()

	// 注册http
	srv := &http.Server{
		Addr:         conf.ApplicationConfig.Host + ":" + conf.ApplicationConfig.Port,
		Handler:      r,
		ReadTimeout:  time.Duration(conf.ApplicationConfig.ReadTimeout) * time.Second,
		WriteTimeout: time.Duration(conf.ApplicationConfig.WriterTimeout) * time.Second,
	}

	go func() {
		//服务器连接
		if err := srv.ListenAndServe(); err != nil && !errors.Is(err, http.ErrServerClosed) {
			panic(fmt.Sprintf("Listen: %s\n", err))
		}
	}()

	fmt.Printf("%s Server Run http://%s:%s/ \r\n",
		time.Now().String(),
		conf.ApplicationConfig.Host,
		conf.ApplicationConfig.Port)
	fmt.Printf("%s Swagger URL http://%s:%s/swagger/index.html \r\n",
		time.Now().String(),
		conf.ApplicationConfig.Host,
		conf.ApplicationConfig.Port)
	fmt.Printf("%s Enter Control + C Shutdown Server \r\n", time.Now().String())

	// 加载机器人启动
	go func() {
		TgBotSetup()
	}()

	// 等待中断信号可以优雅地关闭服务器
	quit := make(chan os.Signal)
	signal.Notify(quit, os.Interrupt)
	<-quit
	fmt.Printf("%s 服务器正在关闭......", time.Now().String())

	ctx, cancel := context.WithTimeout(context.Background(), 5*time.Second)
	defer cancel()
	if err := srv.Shutdown(ctx); err != nil {
		panic(fmt.Sprintf("服务器关闭:%v", err))
	}
	fmt.Printf("%s 服务器已关闭.", time.Now().String())
	return nil
}
