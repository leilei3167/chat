package db

import (
	"github.com/jinzhu/gorm"
	_ "github.com/jinzhu/gorm/dialects/mysql"
	"github.com/sirupsen/logrus"
	"sync"
	"time"
)

var dbMap = map[string]*gorm.DB{}
var dbLock sync.Mutex

func init() { //执行初始化数据库
	initDB("gochat")
}
func initDB(dbName string) {
	var err error
	dbLock.Lock()
	defer dbLock.Unlock()
	dbMap[dbName], err = gorm.Open("mysql", "root:8888@tcp(127.0.0.1:3306)/gochat?charset=utf8mb4&parseTime=True&loc=Local")
	if err != nil {
		logrus.Fatal("初始化db错误:", err)
	}
	//配置数据库
	dbMap[dbName].DB().SetMaxIdleConns(4)
	dbMap[dbName].DB().SetMaxOpenConns(20)
	dbMap[dbName].DB().SetConnMaxLifetime(time.Second * 8)

	dbMap[dbName].LogMode(true) //暂时开启详细日志
}

func GetDb(dbName string) (db *gorm.DB) {
	if db, ok := dbMap[dbName]; ok {
		return db
	} else {
		return nil
	}
}

type DbGoChat struct {
}

func (*DbGoChat) GetDbName() string {
	return "gochat"
}
