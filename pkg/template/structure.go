package template

type Column struct {
	ColumnName    string `json:"column_name"`
	DataType      string `json:"data_type"`
	ColumnKey     string `json:"column_key"`
	ColumnComment string `json:"column_comment"`
}

type StructField struct {
	GoFieldName string
	GoFieldType string
	GormTag     string
	JsonTag     string
	Comment     string
}

type StructTemplateData struct {
	PackageName       string
	ModelName         string
	TableName         string
	Columns           []StructField
	AdminListRequest  string
	AdminListResponse string
}

type RouterTemplateData struct {
	PackageName       string
	RouterMark        string
	StructRouterName  string
	BaseDomain        string
	ApiCreateFuncName string
	ApiListFuncName   string
	ApiOneFuncName    string
	ApiDelFuncName    string
	ApiUpdateFuncName string
}

type ApiTemplateData struct {
	PackageName       string
	BaseModelName     string
	ModelName         string
	OneFuncName       string
	ListFuncName      string
	DelFuncName       string
	UpdateFuncName    string
	CreateFuncName    string
	AdminListResponse string
	AdminListRequest  string
}

type ResponseTemplateData struct {
	PackageName      string
	ResponseName     string
	ResponseDataName string
	Columns          []StructField
}

type RequestTemplateData struct {
	RequestDataName string
}
