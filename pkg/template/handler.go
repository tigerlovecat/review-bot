package template

import (
	"bufio"
	"fmt"
	"log"
	"os"
	"path/filepath"
	"strings"
	"text/template"
)

type HandlerTemplate struct {

	// 项目前缀 + 名称 + 位置 + 模板
	ProjectPrefix       string
	ProjectNameNoPrefix string
	TableName           string
	ProjectPath         string
	TemplateSamples     string
	PathBase            string
	FileName            string
	RouterBase          string

	// Model 相关
	ModelFuncRemark    string
	ModelFuncName      string
	ModelStructColumns []Column
	ModelRequestName   string
	ModelResponseName  string

	// Router 相关
	RouterFuncName   string
	RouterBaseDomain string
	RouterInsertInit string

	// Api 相关
	ApiBaseModelParam string
	ApiBaseModelName  string
	ApiListFuncName   string
	ApiOneFuncName    string
	ApiDelFuncName    string
	ApiUpdateFuncName string
	ApiCreateFuncName string
	ApiListRequest    string
	ApiListResponse   string

	// Response 相关
	ResponseFuncName    string
	ResponseFuncColumns []Column
	ResponseDataName    string

	// Request 相关
	RequestFuncName string
}

// NewHandlerTemplate 初始化
func NewHandlerTemplate(projectPrefix, projectNameNoPrefix, tableName string) *HandlerTemplate {

	// eg: 输出 AstroPoi
	tableName1 := toUpperCamelCase(convertToCamelCase(tableName[2:]))
	// eg: 输出 astroPoi
	tableName2 := convertToCamelCase(tableName[2:])
	// eg: 输出 XAstroPoi
	tableName3 := strings.Title(convertToCamelCase(tableName))
	// 数据库数据
	columns, tableComment, err := getTableColumns(tableName)
	if err != nil {
		fmt.Printf("数据库数据 err: {%v}", err.Error())
		return nil
	}
	return &HandlerTemplate{
		ProjectPrefix:       projectPrefix,
		ProjectNameNoPrefix: projectNameNoPrefix,
		TableName:           tableName,
		ProjectPath:         "TODO",
		FileName:            tableName2 + ".go",
		RouterBase:          tableName1,

		ModelFuncRemark:    tableComment + "相关",
		ModelFuncName:      tableName3 + "Model",
		ModelStructColumns: columns,
		ModelRequestName:   "Admin" + tableName1 + "ListRequest",
		ModelResponseName:  tableName1 + "Response",

		RouterFuncName:   tableName1 + "Router",
		RouterBaseDomain: tableName2,
		RouterInsertInit: tableName1 + "Router(g)",

		ApiBaseModelParam: "base" + tableName3 + "Model",
		ApiBaseModelName:  tableName3 + "Model",
		ApiListFuncName:   "List" + tableName1,
		ApiOneFuncName:    "One" + tableName1,
		ApiDelFuncName:    "Del" + tableName1,
		ApiUpdateFuncName: "Update" + tableName1,
		ApiCreateFuncName: "Create" + tableName1,
		ApiListRequest:    "Admin" + tableName1 + "ListRequest",
		ApiListResponse:   tableName1 + "Response",

		ResponseFuncName:    "Admin" + tableName1 + "Response",
		ResponseFuncColumns: columns,
		ResponseDataName:    tableName1 + "Response",

		RequestFuncName: "Admin" + tableName1 + "ListRequest",
	}
}

// PrepareTemplate 准备模板
func (h *HandlerTemplate) PrepareTemplate() (tmpl *template.Template, file *os.File, err error) {
	tmpl, err = createTemplate("api", h.TemplateSamples)
	if err != nil {
		return
	}
	// 蛇形转驼峰
	distFileName := filepath.Join(h.PathBase, h.FileName)
	file, err = checkAndCreateFile(distFileName)
	if err != nil || file == nil {
		return
	}
	return
}

// GenerateApiCode 生成对应api层的代码
func (h *HandlerTemplate) GenerateApiCode() error {
	h.PathBase = "./api/"
	h.TemplateSamples = apiTemplate
	tmpl, file, err := h.PrepareTemplate()
	if err != nil {
		handleErr(err, "prepareData 数据准备阶段出现异常")
		return err
	}
	defer file.Close()
	// ... 数据处理和模板执行的代码 ...
	data := ApiTemplateData{
		PackageName:       h.ProjectNameNoPrefix,
		BaseModelName:     h.ApiBaseModelParam,
		ModelName:         h.ApiBaseModelName,
		OneFuncName:       h.ApiOneFuncName,
		ListFuncName:      h.ApiListFuncName,
		DelFuncName:       h.ApiDelFuncName,
		UpdateFuncName:    h.ApiUpdateFuncName,
		CreateFuncName:    h.ApiCreateFuncName,
		AdminListResponse: h.ApiListResponse,
		AdminListRequest:  h.ApiListRequest,
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		handleErr(err, "tmpl.Execute err!")
		return err
	}
	return nil
}

// GenerateResponseCode 生成对应的response层的代码
func (h *HandlerTemplate) GenerateResponseCode() error {
	h.PathBase = "./pkg/app/response/"
	h.TemplateSamples = responseTemplate
	tmpl, file, err := h.PrepareTemplate()
	if err != nil {
		handleErr(err, "prepareData 数据准备阶段出现异常")
		return err
	}
	defer file.Close()
	var fields []StructField
	for _, col := range h.ModelStructColumns {
		if col.ColumnName == "created_at" || col.ColumnName == "updated_at" {
			continue
		}
		field := StructField{
			GoFieldName: toUpperCamelCase(convertToCamelCase(col.ColumnName)),
			GoFieldType: mapSQLTypeToGoType(col.DataType),
			GormTag:     col.ColumnName,
			JsonTag:     col.ColumnName,
			Comment:     col.ColumnComment,
		}
		fields = append(fields, field)
	}
	requestData := ResponseTemplateData{
		ResponseDataName: h.ResponseDataName,
		Columns:          fields,
		PackageName:      "response",
		ResponseName:     h.ResponseFuncName,
	}
	err = tmpl.Execute(file, requestData)
	if err != nil {
		fmt.Printf("tmpl.Execute err: {%v}", err.Error())
		return err
	}
	return nil
}

// GenerateRequestCode 生成对应的request层的代码
func (h *HandlerTemplate) GenerateRequestCode() error {
	h.PathBase = "./pkg/app/request/"
	h.TemplateSamples = requestTemplate
	tmpl, file, err := h.PrepareTemplate()
	if err != nil {
		handleErr(err, "prepareData 数据准备阶段出现异常")
		return err
	}
	defer file.Close()
	requestData := RequestTemplateData{
		RequestDataName: h.RequestFuncName,
	}
	err = tmpl.Execute(file, requestData)
	if err != nil {
		fmt.Printf("tmpl.Execute err: {%v}", err.Error())
		return err
	}
	return nil
}

// GenerateRouterCode 生成对应的router层的代码
func (h *HandlerTemplate) GenerateRouterCode() error {
	h.PathBase = "./router/"
	h.TemplateSamples = routerTemplate
	tmpl, file, err := h.PrepareTemplate()
	if err != nil {
		handleErr(err, "prepareData 数据准备阶段出现异常")
		return err
	}
	defer file.Close()
	data := RouterTemplateData{
		PackageName:       h.ProjectNameNoPrefix,
		RouterMark:        h.ModelFuncRemark,
		StructRouterName:  h.RouterFuncName,
		BaseDomain:        h.RouterBaseDomain,
		ApiDelFuncName:    h.ApiDelFuncName,
		ApiListFuncName:   h.ApiListFuncName,
		ApiOneFuncName:    h.ApiOneFuncName,
		ApiUpdateFuncName: h.ApiUpdateFuncName,
		ApiCreateFuncName: h.ApiCreateFuncName,
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("tmpl.Execute err: {%v}", err.Error())
		return err
	}
	return nil
}

// GenerateStructCode Implement a function to map SQL types to Go types and generate tags
func (h *HandlerTemplate) GenerateStructCode() error {
	h.PathBase = "./internal/model/"
	h.TemplateSamples = modelTemplate
	tmpl, file, err := h.PrepareTemplate()
	if err != nil {
		handleErr(err, "prepareData 数据准备阶段出现异常")
		return err
	}
	defer file.Close()
	var fields []StructField
	for _, col := range h.ModelStructColumns {
		if col.ColumnName == "created_at" || col.ColumnName == "updated_at" {
			continue
		}
		field := StructField{
			GoFieldName: toUpperCamelCase(convertToCamelCase(col.ColumnName)),
			GoFieldType: mapSQLTypeToGoType(col.DataType),
			GormTag:     col.ColumnName,
			JsonTag:     col.ColumnName,
			Comment:     col.ColumnComment,
		}
		fields = append(fields, field)
	}
	data := StructTemplateData{
		PackageName:       h.ProjectNameNoPrefix,
		ModelName:         h.ModelFuncName,
		Columns:           fields,
		TableName:         h.TableName,
		AdminListResponse: h.ModelResponseName,
		AdminListRequest:  h.ModelRequestName,
	}
	err = tmpl.Execute(file, data)
	if err != nil {
		fmt.Printf("tmpl.Execute err: {%v}", err.Error())
		return err
	}
	return nil
}

// GenerateRouterTurnCode 新增路由中转层代码
func (h *HandlerTemplate) GenerateRouterTurnCode() error {
	filePath := "./router/init_router.go" // 替换为你的文件路径
	tempFilePath := filePath + ".tmp"

	file, err := os.Open(filePath)
	if err != nil {
		log.Fatalf("failed to open file: %s", err)
	}
	defer file.Close()

	tempFile, err := os.Create(tempFilePath)
	if err != nil {
		log.Fatalf("failed to create temp file: %s", err)
	}
	defer tempFile.Close()

	scanner := bufio.NewScanner(file)
	writer := bufio.NewWriter(tempFile)

	// 读取文件并在特定位置插入代码
	inserted := false
	for scanner.Scan() {
		line := scanner.Text()
		if _, err := writer.WriteString(line + "\n"); err != nil {
			log.Fatalf("failed to write to temp file: %s", err)
		}

		// 检测到 TODO 注释后插入代码
		if !inserted && strings.Contains(line, "TODO 注册业务路由") {
			codeStr := "\n\t// " + h.ModelFuncRemark
			codeStr += "\n\t" + h.RouterInsertInit + "\n"
			if _, err := writer.WriteString(codeStr); err != nil {
				log.Fatalf("failed to write to temp file: %s", err)
			}
			inserted = true
		}
	}

	if err := scanner.Err(); err != nil {
		log.Fatalf("error while reading file: %s", err)
	}

	if err := writer.Flush(); err != nil {
		log.Fatalf("failed to flush to temp file: %s", err)
	}

	// 替换原文件
	if err := os.Rename(tempFilePath, filePath); err != nil {
		log.Fatalf("failed to replace the original file: %s", err)
	}
	return nil
}
