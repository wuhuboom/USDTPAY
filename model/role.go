package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Role struct {
	ID      int    `json:"id" gorm:"primaryKey;comment:'主键'"`
	Name    string `json:"name" gorm:"uniqueIndex"`
	Status  int    `json:"status" gorm:"default:2"` //状态 1 禁用 2正常
	Created int64
}

func CheckIsExistModelRole(db *gorm.DB) {
	if db.Migrator().HasTable(&Role{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Role{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&Role{})

		db.Save(&Role{ID: 1,
			Name:    "超级管理员",
			Status:  2,
			Created: time.Now().Unix()})
	}
}

func (r *Role) GetName(db *gorm.DB) string {
	db.Where("id=?", r.ID).First(r)
	return r.Name
}
