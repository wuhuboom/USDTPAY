package model

import (
	"fmt"
	"gorm.io/gorm"
)

type Config struct {
	ID           int     `json:"id" gorm:"primaryKey;comment:'主键'"`
	MaxPond      int     `gorm:"default:1000" json:"maxPond"`   //通用池子的大小  默认值 1000
	Expiration   int64   `gorm:"default:30" json:"expiration"`  //通用池子的订单过期时间  30 分钟
	PondAmount   float64 `gorm:"default:5" json:"pondAmount"`   //池子的金额分界点  5U
	GoogleSwitch int     `gorm:"default:1" json:"googleSwitch"` //谷歌开光 默认1关闭 2开启
}

// CheckIsExistModelConfig CheckIsExistModelAdmin 创建
func CheckIsExistModelConfig(db *gorm.DB) {
	if db.Migrator().HasTable(&Config{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Config{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&Config{})
		db.Save(&Config{ID: 1})
	}
}

func (c *Config) GetConfig(db *gorm.DB) *Config {
	c.MaxPond = 1000
	c.PondAmount = 10
	c.Expiration = 60
	db.Where("id=?", 1).First(c)
	return c
}
