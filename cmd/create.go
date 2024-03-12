package cmd

import (
	"fmt"
	"github.com/spf13/cobra"
	"log"
	"os"
	cfg "review-bot/config"
	"review-bot/database"
	"review-bot/pkg/template"
)

var (
	CreateCmd = &cobra.Command{
		Use:     "create",
		Short:   "脚手架生成代码",
		Example: "./wb create x_astro_poi",
		PreRun: func(cmd *cobra.Command, args []string) {
			usageStr := `starting create server code`
			log.Printf("%s\n", usageStr)
			// 初始化配置文件
			cfg.Setup(config)
			// 数据初始化
			database.MysqlSetup()
		},
		RunE: func(cmd *cobra.Command, args []string) error {
			return run()
		},
	}
)

func run() error {
	fmt.Println("脚手架生成代码...后面的参数为：")
	if len(os.Args) >= 3 && os.Args[1] == "create" {
		for i := 2; i < len(os.Args); i++ {
			createCode(os.Args[i])
		}
	} else {
		fmt.Println("没有提供第二个参数。")
	}
	return nil
}

// 自动生成代码核心库
func createCode(tableName string) error {

	fmt.Printf("正在准备生成 table_name = %v 的相关代码 \n", tableName)
	handler := template.NewHandlerTemplate(cfg.ApplicationConfig.Prefix, cfg.ApplicationConfig.Name, tableName)

	err := handler.GenerateRequestCode()
	if err != nil {
		fmt.Printf("\n[自动生成代码] -- 插入pkg/app/request类err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入pkg/app/request类完成 \n")

	err = handler.GenerateResponseCode()
	if err != nil {
		fmt.Printf("[自动生成代码] -- 插入pkg/app/response类 err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入pkg/app/response类完成 \n")

	err = handler.GenerateStructCode()
	if err != nil {
		fmt.Printf("[自动生成代码] -- 插入internal/model类 err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入internal/model类完成 \n")

	err = handler.GenerateApiCode()
	if err != nil {
		fmt.Printf("[自动生成代码] -- 插入api类 err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入api类完成 \n")

	err = handler.GenerateRouterCode()
	if err != nil {
		fmt.Printf("[自动生成代码] -- 插入router类 err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入router类完成 \n")

	err = handler.GenerateRouterTurnCode()
	if err != nil {
		fmt.Printf("[自动生成代码] -- 插入router中转层代码 err: {%v} \n", err.Error())
		return err
	}
	fmt.Printf("\n[自动生成代码] -- 插入router中转层代码完成 \n")

	return nil
}
