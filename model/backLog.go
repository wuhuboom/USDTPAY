package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type BackLog struct {
	ID          int    `gorm:"primaryKey;comment:'主键'" json:"id"`
	BackContent string `json:"back_content" gorm:"type:text"`  // 发包内容
	JsonContent string ` gorm:"type:text" json:"json_content"` //解析内容
	Kinds       int    `json:"kinds"`                          //1回调数据  2清零日志
	TxHash      string `json:"tx_hash" gorm:"uniqueIndex"`     // 回调的Hash值
	Created     int64  `json:"created"`
}

func CheckIsExistModelBackLog(db *gorm.DB) {
	if db.Migrator().HasTable(&BackLog{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&BackLog{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&BackLog{})
	}
}

// CreatedBackLog 创建
func (receiver *BackLog) CreatedBackLog(db *gorm.DB) {
	receiver.Created = time.Now().Unix()
	db.Save(receiver)
}
