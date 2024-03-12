package model

import (
	"fmt"
	"reflect"
	"review-bot/database"
	"time"
)

type BaseModel struct {
	CreatedAt        time.Time `gorm:"column:created_at" json:"created_at" example:"2022-12-30 12:23:23"` //创建时间
	UpdatedAt        time.Time `gorm:"column:updated_at" json:"updated_at" example:"2022-12-30 12:23:23"` //修改时间
	CreatedAtFroShow string    `json:"created_at_fro_show" gorm:"-"`
	UpdatedAtFroShow string    `json:"updated_at_fro_show" gorm:"-"`
}

// CommonInterface 通用的DAO接口，不再与具体的表绑定
type CommonInterface interface {
	FindAll() error
	FindById() error
	Insert() error
	Update() error
	Delete() error
}

// CommonDAO 通用的DAO实现，通过传入结构体指针来决定操作哪张数据表
type CommonDAO struct {
	TableName string
	ID        int64
	Base      BaseModel
	Data      interface{}
}

func (dao *CommonDAO) FindAll() error {
	// 根据传递过来的表名参数来确定查询哪张表
	typeTmp := reflect.ValueOf(dao.Data)
	data := reflect.New(typeTmp.Type()).Elem()
	err := database.Eloquent.Table(dao.TableName).Find(data.Addr().Interface()).Error
	if err != nil {
		return err
	}
	dao.Data = data.Interface()
	return nil
}

func (dao *CommonDAO) FindById() error {
	// 根据传递过来的结构体名称来确定查询哪张表，并将查询结果转化为interface{}类型返回
	typeTmp := reflect.ValueOf(dao.Data)
	data := reflect.New(typeTmp.Type()).Elem()
	err := database.Eloquent.Table(dao.TableName).Where("id", dao.ID).First(data.Addr().Interface()).Error
	dao.Data = data.Interface()
	return err
}

func (dao *CommonDAO) Insert() error {
	// 根据传递过来的结构体名称来确定插入哪张表，并将对应传递过来的数据插入到数据库中
	mapValue := reflect.ValueOf(dao.Data)
	if mapValue.Kind() != reflect.Map {
		return fmt.Errorf("invalid data type: %T", dao.Data)
	}

	createData := make(map[string]interface{})
	for _, key := range mapValue.MapKeys() {
		createData[key.String()] = mapValue.MapIndex(key).Interface()
	}

	// 使用db.Create(obj) 插入数据
	result := database.Eloquent.Table(dao.TableName).Create(createData).Order("id desc").Take(&createData)
	if result.Error != nil {
		return result.Error
	}

	// 获取自增长ID
	dao.Data = createData

	return nil
}

func (dao *CommonDAO) Update() error {
	// 根据传递过来的结构体及名称来确定更新哪张表中的数据
	return database.Eloquent.Table(dao.TableName).Debug().Where("id", dao.ID).Updates(dao.Data).Error
}

func (dao *CommonDAO) Delete() error {
	// 根据传递过来的结构体名称来确定删除哪张表中的数据
	return database.Eloquent.Table(dao.TableName).Where("id", dao.ID).Delete(dao.Data).Error
}
