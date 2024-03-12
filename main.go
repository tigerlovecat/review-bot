package main

import (
	"errors"
	"fmt"
	"log"
	"os"
	"review-bot/cmd"
	"time"

	"github.com/spf13/cobra"
)

var (
	// 启动
	rootCmd = &cobra.Command{
		Use:               "wb",
		Short:             "-v",
		SilenceUsage:      true,
		DisableAutoGenTag: true,
		Long:              `wb`,
		Args: func(cmd *cobra.Command, args []string) error {
			if len(args) < 1 {
				return errors.New("至少需要一个参数")
			} else {
				return errors.New("参数错误")
			}
		},
		PersistentPostRunE: func(cmd *cobra.Command, args []string) error {
			return nil
		},
		Run: func(cmd *cobra.Command, args []string) {
			usageStr := `欢迎使用fanxing, 可以使用 -h 查看命令`
			fmt.Println(usageStr)
		},
	}
)

func init() {
	fmt.Println("默认走init 函数，启动项目")
	rootCmd.AddCommand(cmd.StartCmd)
	rootCmd.AddCommand(cmd.CreateCmd)
}

func main() {
	defer func() {
		if p := recover(); p != nil {
			// 处理panic信息
			log.Printf("panic recover! p: %v", p)
		}
	}()
	location, err := time.LoadLocation("Asia/Shanghai")
	if err != nil {
		location = time.FixedZone("CST", 8*3600)
	}
	time.Local = location
	Execute()
}

func Execute() {
	if err := rootCmd.Execute(); err != nil {
		os.Exit(-1)
	}
}
