package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Log struct {
	ID      int    `gorm:"primaryKey;comment:'主键'" json:"id"`
	Content string `json:"content" gorm:"type:text"` //内容
	Ips     string `json:"ips"`
	Kinds   int    `json:"kinds"` //日志种类  1登录日志 2系统错误日志 3获取地址日志  4资金归集日志
	Created int64  `json:"created"`
}

func CheckIsExistModelLog(db *gorm.DB) {
	if db.Migrator().HasTable(&Log{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Log{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&Log{})

	}
}

func (receiver *Log) CreatedLogs(db *gorm.DB) {
	receiver.Created = time.Now().Unix()
	db.Save(receiver)
}
