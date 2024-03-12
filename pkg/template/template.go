package template

// model模板
const modelTemplate = `package model

import (
	"{{.PackageName}}/database"
	"{{.PackageName}}/pkg/app/request"
	"{{.PackageName}}/pkg/app/response"
)

type {{.ModelName}} struct {
	BaseModel
{{- range .Columns}}
    {{.GoFieldName}} {{.GoFieldType}} ` + "`" + `gorm:"{{.GormTag}}" json:"{{.JsonTag}}"` + "`" + ` // {{.Comment}}
{{- end}}
}

func (c *{{.ModelName}}) TableName() string {
	return "{{.TableName}}"
}

func (c *{{.ModelName}}) List(params request.{{.AdminListRequest}}) (result []response.{{.AdminListResponse}}, total int64, err error) {
	conn := database.Eloquent.Table(c.TableName())
	offset := (params.Page - 1) * params.PageSize
	err = conn.Offset(offset).Limit(params.PageSize).Scan(&result).Offset(-1).Limit(-1).Count(&total).Error
	if err != nil {
		return result, total, err
	}
	return result, total, nil
}
`

// 新增路由层模板
const routerTemplate = ` package router

import (
	"{{.PackageName}}/api"
	"{{.PackageName}}/internal/middlewares"
	"github.com/gin-gonic/gin"
)

func {{.StructRouterName}}(r *gin.RouterGroup) {
	withAuthHandler := r.Group("/{{.BaseDomain}}", middlewares.CheckToken())
	{
		withAuthHandler.GET("list", api.{{.ApiListFuncName}})
		withAuthHandler.GET("one", api.{{.ApiOneFuncName}})
		withAuthHandler.POST("create", api.{{.ApiCreateFuncName}})
		withAuthHandler.POST("update", api.{{.ApiUpdateFuncName}})
		withAuthHandler.GET("del", api.{{.ApiDelFuncName}})
	}
}
`

// 新增api层模板
const apiTemplate = ` package api

import (
	"{{.PackageName}}/internal/model"
	"{{.PackageName}}/pkg/app"
	"{{.PackageName}}/pkg/app/request"
	"{{.PackageName}}/pkg/app/response"
	"{{.PackageName}}/pkg/app/utils"
	"fmt"
	"github.com/gin-gonic/gin"
	"math"
	"strconv"
)

var (
	{{.BaseModelName}} model.{{.ModelName}}
)

// {{.ListFuncName}} 列表
func {{.ListFuncName}}(c *gin.Context) {
	var params request.{{.AdminListRequest}}
	err := c.ShouldBind(&params)
	if err != nil {
		app.XError(c, "ParamsHandlerErr")
		return
	}
	{{.BaseModelName}} = model.{{.ModelName}}{}
	result, total, err := {{.BaseModelName}}.List(params)
	if err != nil {
		app.XError(c, "DBGetErr")
		return
	}
	if total == 0 {
		result = []response.{{.AdminListResponse}}{}
	}
	app.XOK(c, &response.ListPageResponse{
		Total:     total,
		CurrentPage: params.Page,
		TotalPage: int64(math.Ceil(float64(total) / float64(params.PageSize))),
		List:      result,
	}, "")
	return
}

// {{.OneFuncName}} 单条信息
func {{.OneFuncName}}(c *gin.Context) {
	{{.BaseModelName}} = model.{{.ModelName}}{}
	idStr := c.Query("id")
	idInt64, _ := strconv.ParseInt(idStr, 10, 64)
	handler := new(model.CommonDAO)
	handler.TableName = {{.BaseModelName}}.TableName()
	handler.Data = {{.BaseModelName}}
	handler.ID = idInt64
	err := handler.FindById()
	if err != nil {
		app.XError(c, "DBGetErr")
		return
	}
	data, ok := handler.Data.(model.{{.ModelName}})
	if !ok {
		app.XError(c, "DataUnmarshalErr")
		return
	}
	data.CreatedAtFroShow = data.CreatedAt.Format(utils.BaseFormatTime)
	data.UpdatedAtFroShow = data.UpdatedAt.Format(utils.BaseFormatTime)
	fmt.Printf("data: {%+v}", data)
	app.XOK(c, &data, "")
	return
}

// {{.DelFuncName}} 删除
func {{.DelFuncName}}(c *gin.Context) {
	{{.BaseModelName}} = model.{{.ModelName}}{}
	idStr := c.Query("id")
	idInt64, _ := strconv.ParseInt(idStr, 10, 64)
	handler := new(model.CommonDAO)
	handler.TableName = {{.BaseModelName}}.TableName()
	handler.Data = model.{{.ModelName}}{}
	handler.ID = idInt64
	err := handler.Delete()
	if err != nil {
		app.XError(c, "DBGetErr")
		return
	}
	app.XOK(c, nil, "")
	return
}

// {{.UpdateFuncName}} 更新
func {{.UpdateFuncName}}(c *gin.Context) {
	{{.BaseModelName}} = model.{{.ModelName}}{}
	var updateMap map[string]interface{}
	err := c.ShouldBindJSON(&updateMap)
	if err != nil {
		app.XError(c, "ParamsHandlerErr")
		return
	}
	handler := new(model.CommonDAO)
	handler.TableName = {{.BaseModelName}}.TableName()
	if _, ok := updateMap["id"]; !ok {
		app.XError(c, "ParamsNoIDErr")
		return
	}
	handler.ID = int64(updateMap["id"].(float64))
	delete(updateMap, "id")
	handler.Data = updateMap
	err = handler.Update()
	if err != nil {
		app.XError(c, "DBGetErr")
		return
	}
	app.XOK(c, &handler.Data, "")
	return
}

// {{.CreateFuncName}} 新建
func {{.CreateFuncName}}(c *gin.Context) {
	{{.BaseModelName}} = model.{{.ModelName}}{}
	var createMap map[string]interface{}
	err := c.ShouldBindJSON(&createMap)
	if err != nil {
		app.XError(c, "ParamsHandlerErr")
		return
	}
	handler := new(model.CommonDAO)
	handler.TableName = {{.BaseModelName}}.TableName()
	handler.Data = createMap
	err = handler.Insert()
	if err != nil {
		app.XError(c, "DBGetErr")
		return
	}
	app.XOK(c, &handler.Data, "")
	return
}
`

// 新增response层模板
const responseTemplate = ` package response

type {{.ResponseDataName}} struct {
{{- range .Columns}}
    {{.GoFieldName}} {{.GoFieldType}} ` + "`" + `gorm:"{{.GormTag}}" json:"{{.JsonTag}}"` + "`" + ` // {{.Comment}}
{{- end}}
}

`

// 新增request层模板
const requestTemplate = ` package request

type {{.RequestDataName}} struct {
	Page int ` + "`" + `form:"page" ` + "`" + `
	PageSize int ` + "`" + `form:"page_size" ` + "`" + `
}
`
