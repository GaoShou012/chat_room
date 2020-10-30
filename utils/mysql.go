package utils

import (
	_ "github.com/go-sql-driver/mysql"
	"github.com/jinzhu/gorm"
)

var (
	DB *gorm.DB
)

//func InitMysql(user string, password string, addr string, port int, database string) {
//	var err error
//	dns := fmt.Sprintf("%s:%s@tcp(%s:%d)/%s?charset=utf8mb4&parseTime=True&loc=Local", user, password, addr, port, database)
//	db, err := gorm.Open("mysql", dns)
//	if err != nil {
//		glog.Error("failed to conn mysql, %v", err)
//		return
//	}
//	db.DB().SetMaxIdleConns(10)
//	db.DB().SetMaxOpenConns(100)
//	db.LogMode(true)
//	DB = db
//}
