package model

import (
	"example.com/m/tools"
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Admin struct {
	ID         int    `json:"id" gorm:"primaryKey;comment:'主键'"`
	Token      string `json:"token"`
	Username   string `json:"username" gorm:"uniqueIndex"`
	Password   string `json:"password"`
	Status     int    `gorm:"default:1" json:"status"` //状态 1正常 2禁用
	GoogleCode string `json:"google_code"`
	RoleId     int    `json:"role_id"` //角色Id
	Created    int64  `json:"created"`
	Updated    int64  `json:"updated"`
	RoleName   string `gorm:"-" json:"role_name"`
}

// CheckIsExistModelAdmin 创建
func CheckIsExistModelAdmin(db *gorm.DB) {
	if db.Migrator().HasTable(&Admin{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Admin{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&Admin{})
		db.Save(&Admin{
			Username: "ace001",
			Password: tools.MD5("ace001"),
			Token:    string(tools.RandString(36)),
			Created:  time.Now().Unix(), RoleId: 1})
	}
}
