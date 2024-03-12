package database

import (
	"fmt"
	"log"
	"os"
	"strings"
	"time"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
	"gorm.io/gorm/callbacks"
	gormlog "gorm.io/gorm/logger"
	"gorm.io/plugin/dbresolver"
	"review-bot/config"
)

func (e *Mysql) Setup() {

	var (
		err error
		db  Database
	)

	db = new(Mysql)
	MysqlConn = db.GetMasterConnect()
	Eloquent, err = db.Open(MysqlConn)
	newLogger := gormlog.New(
		log.New(os.Stdout, "\r\n", log.LstdFlags), // io writer
		gormlog.Config{
			SlowThreshold:             time.Second,  // 慢 SQL 阈值
			LogLevel:                  gormlog.Info, // Log level
			Colorful:                  true,         // 禁用彩色打印
			IgnoreRecordNotFoundError: false,
		},
	)
	Eloquent, err = gorm.Open(mysql.Open(MysqlConn), &gorm.Config{
		DisableNestedTransaction: true,
		Logger:                   newLogger,
		PrepareStmt:              false,
	})

	// 开启mysql 语句输出
	Eloquent.Debug()
	// 开启 查询超时标识
	// Eloquent.Callback().Query().Replace("gorm:query", maxExecutionTime)
	// 配置读写分离
	Eloquent.Use(dbresolver.Register(dbresolver.Config{
		Replicas: []gorm.Dialector{mysql.Open(db.GetSlaveConnect())},
		Sources:  []gorm.Dialector{mysql.Open(db.GetMasterConnect())},
	}).
		SetMaxIdleConns(10).
		SetMaxOpenConns(600).
		SetConnMaxLifetime(time.Second * 30))
	// SetMaxIdleConns 最大闲置连接数
	// SetMaxOpenConns 数据库最大连接数
	// SetConnMaxLifetime 数据库链接超时时间
	if err != nil {
		panic(fmt.Sprintf("mysql connect error %w", err))
	} else {
		fmt.Println("mysql connect success!")
	}

	if Eloquent.Error != nil {
		panic(fmt.Sprintf("database error %w", Eloquent.Error))
	}
}

type Mysql struct {
}

func (e *Mysql) Open(conn string) (db *gorm.DB, err error) {
	return gorm.Open(mysql.Open(conn))
}

func (e *Mysql) GetMasterConnect() string {
	return config.DatabaseConfig.Master
}

func (e *Mysql) GetSlaveConnect() string {
	return config.DatabaseConfig.Slave
}

func maxExecutionTime(db *gorm.DB) {

	if db.Error == nil {
		callbacks.BuildQuerySQL(db)
		sql := db.Statement.SQL.String()
		sql = strings.Replace(sql, "SELECT", "SELECT /*+ MAX_EXECUTION_TIME(5000) */", 1)
		if !db.DryRun && db.Error == nil {
			rows, err := db.Statement.ConnPool.QueryContext(db.Statement.Context, sql, db.Statement.Vars...)
			if err != nil {
				db.AddError(err)
				return
			}
			defer rows.Close()
			gorm.Scan(rows, db, 0)
		}
	}
}
