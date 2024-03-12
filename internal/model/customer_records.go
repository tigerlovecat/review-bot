package model

import "review-bot/database"

type CustomerRecordsModel struct {
	Id          int64  `gorm:"id" json:"id"`
	Username    string `gorm:"username" json:"username"`
	Message     string `gorm:"message" json:"message"`
	Reply       string `gorm:"reply" json:"reply"`
	MessageTime int64  `gorm:"message_time" json:"message_time"`
}

func (c *CustomerRecordsModel) TableName() string {
	return "customer_records"
}

func (c *CustomerRecordsModel) List() (list []CustomerRecordsModel, err error) {
	err = database.Eloquent.Table(c.TableName()).Order("id desc").Scan(&list).Error
	return list, err
}

func (c *CustomerRecordsModel) AddRecord() error {
	return database.Eloquent.Table(c.TableName()).Create(&c).Error
}
