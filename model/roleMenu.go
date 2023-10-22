package model

import (
	"fmt"
	"gorm.io/gorm"
)

type RoleMenu struct {
	ID       int    `json:"id" gorm:"primaryKey;comment:'主键'"`
	RoleId   int    `json:"role_id"`
	MenuId   int    `json:"menu_id"`
	MenuName string `json:"menu_name" gorm:"-"`
	Created  int64
}

func CheckIsExistModelRoleMenu(db *gorm.DB) {
	if db.Migrator().HasTable(&RoleMenu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&RoleMenu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&RoleMenu{})

		//初始化超级管理员的
		me := make([]Menu, 0)
		db.Find(&me)
		for _, menu := range me {
			db.Save(&RoleMenu{RoleId: 1, MenuId: int(menu.ID)})
		}

	}
}
