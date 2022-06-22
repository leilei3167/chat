// Package dao 定义logic层对user信息的关系型数据库操作
package dao

import (
	"errors"
	"github.com/jinzhu/gorm"
	"github.com/leilei3167/chat/internal/db"
)

var dbIns = db.GetDb("gochat") //包内使用的数据库连接实例

func init() {
	var u User
	dbIns.AutoMigrate(&u) //初始化表格
}

type User struct {
	UserName string
	Password string
	gorm.Model
	db.DbGoChat
}

func (u *User) TableName() string { //返回表名字,便于构建gorm查询语句
	return "user"
}

// Add 将一个User加入到数据库中
func (u *User) Add() (id int, err error) {
	if u.UserName == "" || u.Password == "" {
		return 0, errors.New("username or password empty")
	}
	//检查是否存在相同用户名,存在则返回其id
	oUser := u.CheckHaveUserName(u.UserName)
	if oUser.ID > 0 {
		return int(oUser.ID), nil
	}
	//否则存入(Table指定表格)
	if err = dbIns.Table(u.TableName()).Create(&u).Error; err != nil {
		return 0, err
	}
	return int(u.ID), nil
}

// CheckHaveUserName 根据用户名查找
func (u *User) CheckHaveUserName(userName string) (data User) {
	dbIns.Table(u.TableName()).Where("user_name=?", userName).First(&data)
	return
}

func (u *User) GetUserNameByUserID(id int) (userName string) {
	var data User
	dbIns.Table(u.TableName()).Where("user_id=?", id).First(&data)
	return data.UserName
}
