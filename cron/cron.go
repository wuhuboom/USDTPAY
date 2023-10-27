package cron

import (
	"example.com/m/model"
	"fmt"
	"github.com/robfig/cron/v3"
	"gorm.io/gorm"
	"time"
)

func CornInit(db *gorm.DB) {
	c := cron.New()
	// 添加定时任务
	_, err := c.AddFunc("0 0 6 * *", func() {
		// 在每天6点钟触发的任务逻辑
		//var Data model.ConsoleManagementData
		//ReturnConsoleManagementData(Data, db)
		var Data model.ConsoleManagementData
		Data.Date = time.Now().AddDate(0, 0, -1).Format("2006-01-02")
		ReturnConsoleManagementData(&Data, db, time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
		Data.CreatedConsoleManagementData(db)
	})
	//var Data model.ConsoleManagementData
	//ReturnConsoleManagementData(&Data, db, time.Now().AddDate(0, 0, -1).Format("2006-01-02"))
	//Data.CreatedConsoleManagementData(db)
	if err != nil {
		fmt.Println("添加定时任务失败：", err)
		return
	}
	// 启动定时器
	c.Start()
	// 程序继续运行，等待定时任务执行
	select {}
}

func ReturnConsoleManagementData(Data *model.ConsoleManagementData, db *gorm.DB, date string) *model.ConsoleManagementData {
	db.Model(&model.PrepaidPhoneOrders{}).
		Where("status =? and date =?", 2, date).
		Count(&Data.TodayPullOrderCountAndSuccess)
	//今日拉去订单个数
	db.Model(&model.PrepaidPhoneOrders{}).
		Where("date =?", date).
		Count(&Data.TodayPullOrderCount)
	//今日拉起订单金额
	db.Table("prepaid_phone_orders").Where("date =?", date).Select("sum(account_orders) as today_pull_order_amount").Scan(&Data)
	//今日成功订单金额
	db.Table("prepaid_phone_orders").Where("date =? and status=?", date, 2).Select("sum(account_practical) as today_pull_order_amount_and_success").Scan(&Data)
	//今日订单支付成功率
	if Data.TodayPullOrderCount == 0 {
		Data.TodaySuccessPer = 0
	} else {
		Data.TodaySuccessPer = float64(Data.TodayPullOrderCountAndSuccess) / float64(Data.TodayPullOrderCount)
	}
	//总成功订单个数
	db.Model(&model.PrepaidPhoneOrders{}).Where("status =? ", 2).Count(&Data.AllPullOrderCountAndSuccess)
	//总订单数
	db.Model(&model.PrepaidPhoneOrders{}).Count(&Data.AllPullOrderCount)
	//总拉起订单金额
	db.Table("prepaid_phone_orders").Select("sum(account_orders) as all_pull_order_amount").Scan(&Data)
	//总成功订单金额
	db.Table("prepaid_phone_orders").Where("status =?", 2).Select("sum(account_practical) as all_pull_order_amount_and_success").Scan(&Data)
	if Data.AllPullOrderCount == 0 {
		Data.AllSuccessPer = 0
	} else {
		Data.AllSuccessPer = float64(Data.AllPullOrderCountAndSuccess) / float64(Data.AllPullOrderCount)
	}
	return Data
}
