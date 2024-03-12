package database

import "gorm.io/gorm"

type Database interface {
	Open(conn string) (db *gorm.DB, err error)
	GetMasterConnect() string
	GetSlaveConnect() string
}
