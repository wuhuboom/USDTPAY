package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type ConsoleManagementData struct {
	ID int `json:"id" gorm:"primaryKey;comment:'主键'"`

	TodayPullOrderCount            int64   `json:"today_pull_order_count"`              //今日拉单总数
	TodayPullOrderCountAndSuccess  int64   `json:"today_pull_order_count_and_success"`  //今日成功支付笔数
	TodayPullOrderAmount           float64 `json:"today_pull_order_amount"`             //今日拉单总金额
	TodayPullOrderAmountAndSuccess float64 `json:"today_pull_order_amount_and_success"` //今日收取金额
	TodaySuccessPer                float64 `json:"today_success_per"`                   //今日订单支付成功率
	AllPullOrderCount              int64   `json:"all_pull_order_count"`                //总拉单总数
	AllPullOrderCountAndSuccess    int64   `json:"all_pull_order_count_and_success"`    //总成功支付笔数
	AllPullOrderAmount             float64 `json:"all_pull_order_amount"`               //总拉单总金额
	AllPullOrderAmountAndSuccess   float64 `json:"all_pull_order_amount_and_success"`   //总收取金额
	AllSuccessPer                  float64 `json:"all_success_per"`                     //总成功率
	Date                           string  `json:"date" gorm:"uniqueIndex"`             //日期
	Created                        int64   `json:"created"`                             //创建时间
}

func CheckIsExistModelConsoleManagementData(db *gorm.DB) {
	if db.Migrator().HasTable(&ConsoleManagementData{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&ConsoleManagementData{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&ConsoleManagementData{})

	}
}

func (c *ConsoleManagementData) CreatedConsoleManagementData(db *gorm.DB) {
	//c.Date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
	c.Created = time.Now().Unix()
	cc := ConsoleManagementData{}
	affected := db.Where("date=?", c.Date).Limit(1).Find(&cc).RowsAffected
	if affected > 0 {
		db.Model(&ConsoleManagementData{}).Where("id=?", cc.ID).Updates(c)
		return
	}
	db.Create(c)
}
