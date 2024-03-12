package template

import (
	"fmt"
	"log"
	"os"
	"review-bot/database"
	"strings"
	"text/template"
	"unicode"
)

// convertToCamelCase 将 snake_case 转换为 CamelCase
func convertToCamelCase(input string) string {
	isToUpper := false
	return strings.Map(func(r rune) rune {
		if r == '_' {
			isToUpper = true
			return -1 // 删除下划线
		}
		if isToUpper {
			isToUpper = false
			return unicode.ToUpper(r)
		}
		return r
	}, input)
}

// toUpperCamelCase 将首字母小写的 CamelCase（小驼峰命名）
func toUpperCamelCase(input string) string {
	if input == "" {
		return ""
	}
	rns := []rune(input)
	rns[0] = unicode.ToUpper(rns[0])
	return string(rns)
}

func mapSQLTypeToGoType(sqlType string) string {
	switch sqlType {
	case "varchar", "text":
		return "string"
	case "int":
		return "int"
	case "tinyint":
		return "int"
	case "bigint":
		return "int64"
	case "boolean":
		return "bool"
	case "datetime", "timestamp":
		return "string"
	// 添加更多映射关系
	default:
		return "interface{}" // 默认类型
	}
}

// handleErr 是一个通用的错误处理函数
func handleErr(err error, message string) error {
	if err != nil {
		log.Printf("%s: {%v}", message, err)
		return err
	}
	return nil
}

// createTemplate 是一个通用的模板创建函数
func createTemplate(name, tmplStr string) (*template.Template, error) {
	tmpl, err := template.New(name).Parse(tmplStr)
	return tmpl, handleErr(err, fmt.Sprintf("template.New %s err", name))
}

// checkAndCreateFile 检查文件是否存在，并在不存在的情况下创建文件
func checkAndCreateFile(filePath string) (*os.File, error) {
	if _, err := os.Stat(filePath); err == nil {
		log.Println("已经存在对应的文件:", filePath)
		return nil, nil
	} else if os.IsNotExist(err) {
		file, err := os.Create(filePath)
		if err != nil {
			return nil, handleErr(err, "os.Create err")
		}
		return file, nil
	} else {
		return nil, handleErr(err, "文件存在与否异常处理")
	}
}

// getTableColumns 从数据库中获取指定表的列信息
func getTableColumns(tableName string) (columns []Column, tableComment string, err error) {
	tableComment = "数据表的描述"
	db := database.Eloquent
	if db == nil {
		return columns, tableComment, fmt.Errorf("database connection is nil")
	}
	// 构建查询SQL语句
	// SELECT COLUMN_NAME, DATA_TYPE, COLUMN_KEY FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = "x_astro";
	query := fmt.Sprintf("SELECT COLUMN_NAME as column_name, DATA_TYPE as data_type, COLUMN_KEY as column_key, COLUMN_COMMENT  as column_comment FROM INFORMATION_SCHEMA.COLUMNS WHERE TABLE_NAME = '%s'", tableName)

	// 使用GORM的Raw和Scan方法执行查询
	result := db.Raw(query).Scan(&columns)
	if result.Error != nil {
		return columns, tableComment, result.Error
	}

	row := db.Raw("SELECT TABLE_COMMENT FROM information_schema.TABLES WHERE TABLE_NAME = ?", tableName).Row()
	row.Scan(&tableComment)
	return columns, tableComment, nil
}
