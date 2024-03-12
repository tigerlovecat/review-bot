package database

func MysqlSetup() {
	//初始化mysql
	var db = new(Mysql)
	db.Setup()
}
