package ioc

import (
	"webook/config"
	"webook/internal/repository/dao"

	"gorm.io/driver/mysql"
	"gorm.io/gorm"
)

func InitDB() *gorm.DB {
	db, err := gorm.Open(mysql.Open(config.Config.DB.DSN))
	if err != nil {
		panic(err)
	}

	db = db.Debug()

	err = dao.InitTables(db)
	if err != nil {
		panic(err)
	}
	return db
}
