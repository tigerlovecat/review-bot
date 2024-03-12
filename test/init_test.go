package test

import (
	"log"
	cfg "review-bot/config"
	"review-bot/database"
	"review-bot/pkg/template"
	"testing"
)

func InitEnv() {
	cfg.Setup("../config/settings.yaml")
	// database.RedisSetup()
	database.MysqlSetup()
	log.Print("启动单元测试")
}

func TestCreateGoNew(t *testing.T) {
	InitEnv()
	tableName := "x_astro_poi"
	t.Logf("准备生成代码 table: %v", tableName)

	handler := template.NewHandlerTemplate("x", "astro", tableName)

	// 第一步：准备阶段
	// 读取数据库 ->  解析到结构体上
	// 生成对应到Request Response

	t.Logf("第一步：1.1 读取数据库 ->  解析到结构体上")
	err := handler.GenerateRequestCode()
	if err != nil {
		t.Logf("第一步：新增pkg/app/request类 err: {%v}", err.Error())
		return
	}
	t.Logf("第一步：1.2 新增pkg/app/request类")
	err = handler.GenerateResponseCode()
	if err != nil {
		t.Logf("第一步：1.3 新增pkg/app/response类 err: {%v}", err.Error())
		return
	}
	t.Logf("第一步：1.3 新增pkg/app/response类")

	// 第二步：新增internal/model类
	err = handler.GenerateStructCode()
	if err != nil {
		t.Logf("第二步：新增internal/model类 err: {%v}", err.Error())
		return
	}

	// 第三步：新增api层代码
	err = handler.GenerateApiCode()
	if err != nil {
		t.Logf("第三步：新增api层代码 err: {%v}", err.Error())
		return
	}

	// 第四步：新增router层代码
	err = handler.GenerateRouterCode()
	if err != nil {
		t.Logf("第四步：新增router层代码 err: {%v}", err.Error())
		return
	}

	// 第五步：新增router中转层代码
	err = handler.GenerateRouterTurnCode()
	if err != nil {
		t.Logf("第五步：新增router中转层代码 err: {%v}", err.Error())
		return
	}

}
