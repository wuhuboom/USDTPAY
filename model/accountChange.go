package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type AccountChange struct {
	ID             int     `gorm:"primaryKey;comment:'主键'"`
	OriginalAmount float64 `gorm:"type:decimal(10,2)"` // 原始金额
	ChangeAmount   float64 `gorm:"type:decimal(10,2)"` // 原始金额
	NowAmount      float64 `gorm:"type:decimal(10,2)"` //现在的金额
	Kinds          int     // 变化种类  1更新余额   2玩家充值
	//ReceiveAddressId int     //地址id
	Created            int64 //创建时间
	ReceiveAddressName string
}

func CheckIsExistModelAccountChange(db *gorm.DB) {
	if db.Migrator().HasTable(&AccountChange{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&AccountChange{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&AccountChange{})
	}
}

// Add 添加账变记录
func (ac *AccountChange) Add(db *gorm.DB) {
	ac.Created = time.Now().Unix()
	db.Save(&ac)
}
