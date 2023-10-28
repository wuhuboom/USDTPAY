package controller

import (
	"example.com/m/common"
	"example.com/m/cron"
	"example.com/m/dao/mysql"
	"example.com/m/dao/redis"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"time"
)

// Login   管理员登录接口
func Login(c *gin.Context) {
	var lo LoginVerify
	//检查参数
	if err := c.ShouldBind(&lo); err != nil {
		tools.ReturnVerifyErrCode(c, err)
		return
	}
	//判断系统是否开启的谷歌验证
	config := model.Config{}
	err := mysql.DB.Where("id=?", 1).First(&config).Error
	if err != nil {
		tools.ReturnError101(c, "System error. Please contact technical")
		return
	}
	admin := model.Admin{}
	err = mysql.DB.Where("username=?", lo.Username).
		Where("password=?", tools.MD5(lo.Password)).
		First(&admin).Error
	if config.GoogleSwitch == 2 {
		//判断这个用户是否已经绑定了谷歌
		if admin.GoogleCode == "" {
			//没有绑定谷歌  所以要返回谷歌的验证码
			if admin.GoogleCode == "" && lo.GoogleSecret == "" {
				secret, _, qrCodeUrl := tools.InitAuth(admin.Username)
				tools.JsonWrite(c, common.NeedGoogleBind, map[string]string{"codeUrl": qrCodeUrl, "googleSecret": secret}, "Please bind your Google account first")
				return

			} else {
				verifyCode, _ := tools.NewGoogleAuth().VerifyCode(lo.GoogleSecret, lo.GoogleCode)
				if !verifyCode {
					tools.ReturnError101(c, "Google verification failure")
					return
				}
				err := mysql.DB.Model(&model.Admin{}).Where("id=?", admin.ID).Updates(
					model.Admin{GoogleCode: lo.GoogleSecret}).Error
				if err != nil {
					tools.ReturnError101(c, err.Error())
					return
				}
			}
		} else {
			//校验谷歌验证
			verifyCode, _ := tools.NewGoogleAuth().VerifyCode(admin.GoogleCode, lo.GoogleCode)
			if !verifyCode {
				tools.ReturnError101(c, "Google verification failure")
				return
			}
		}
	} else {
		//未开启谷歌
		if err != nil {
			tools.ReturnError101(c, "login fail")
			return
		}
	}
	redis.Rdb.Set(c, "AdminToken_"+admin.Token, admin.Username, 24*time.Hour)
	log := model.Log{Content: fmt.Sprintf("用户:%s,登录成功", admin.Username), Kinds: 1, Ips: c.ClientIP()}
	log.CreatedLogs(mysql.DB)
	tools.ReturnError200Data(c, admin, "success")
	return
}

// 获取菜单

// ConsoleManagement 控制台查看
func ConsoleManagement(c *gin.Context) {
	action := c.PostForm("action")
	//获取数据
	if action == "check" {
		var Data ConsoleManagementData
		//今日成功订单个数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).
			Where("status =? and date =?", 2, time.Now().Format("2006-01-02")).
			Count(&Data.TodayPullOrderCountAndSuccess)
		//今日拉去订单个数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).
			Where("date =?", time.Now().Format("2006-01-02")).
			Count(&Data.TodayPullOrderCount)
		//今日拉起订单金额
		mysql.DB.Table("prepaid_phone_orders").Where("date =?", time.Now().Format("2006-01-02")).Select("sum(account_orders) as today_pull_order_amount").Scan(&Data)
		//今日成功订单金额
		mysql.DB.Table("prepaid_phone_orders").Where("date =? and status=?", time.Now().Format("2006-01-02"), 2).Select("sum(account_practical) as today_pull_order_amount_and_success").Scan(&Data)
		//今日订单支付成功率
		if Data.TodayPullOrderCount == 0 {
			Data.TodaySuccessPer = 0
		} else {
			Data.TodaySuccessPer = float64(Data.TodayPullOrderCountAndSuccess) / float64(Data.TodayPullOrderCount)
		}
		//总成功订单个数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).Where("status =? ", 2).Count(&Data.AllPullOrderCountAndSuccess)
		//总订单数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).Count(&Data.AllPullOrderCount)
		//总拉起订单金额
		mysql.DB.Table("prepaid_phone_orders").Select("sum(account_orders) as all_pull_order_amount").Scan(&Data)
		//总成功订单金额
		mysql.DB.Table("prepaid_phone_orders").Where("status =?", 2).Select("sum(account_practical) as all_pull_order_amount_and_success").Scan(&Data)
		if Data.AllPullOrderCount == 0 {
			Data.AllSuccessPer = 0
		} else {
			Data.AllSuccessPer = float64(Data.AllPullOrderCountAndSuccess) / float64(Data.AllPullOrderCount)
		}
		tools.ReturnError200Data(c, Data, "success")
		return
	}
	//更新指定日期
	if action == "updateDate" {
		date := c.PostForm("date")
		if date == "" {
			tools.ReturnError101(c, "输入正确日期")
			return
		}
		var Data model.ConsoleManagementData
		Data.Date = date
		cron.ReturnConsoleManagementData(&Data, mysql.DB, date)
		Data.CreatedConsoleManagementData(mysql.DB)
		tools.ReturnError200(c, "OK")
		return
	}
	//获取每日数据分页
	if action == "everyday" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.ConsoleManagementData, 0)
		Db := mysql.DB
		var total int64
		Db.Table("console_management_data").Count(&total)
		Db = Db.Model(&model.ConsoleManagementData{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
}

// GetAllMoney 获取总余额
func GetAllMoney(c *gin.Context) {
	rec := make([]model.ReceiveAddress, 0)
	err := mysql.DB.Find(&rec).Error
	if err != nil {
		tools.ReturnError101(c, "获取失败")
		return
	}
	var sumMoney float64
	for _, v := range rec {
		sumMoney = sumMoney + v.Money
	}
	tools.ReturnError200Data(c, sumMoney, "获取成功")
	return
}
