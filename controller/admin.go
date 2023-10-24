package controller

import (
	"example.com/m/common"
	"example.com/m/dao/mysql"
	"example.com/m/dao/redis"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/gin-gonic/gin"
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
			if admin.GoogleCode == "" {
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
			Data.TodaySuccessPer = float64(Data.TodayPullOrderCountAndSuccess / Data.TodayPullOrderCount)
		}
		//总成功订单个数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).Where("status =? ", 2).Count(&Data.AllPullOrderCount)
		//总订单数
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).Count(&Data.AllPullOrderCount)
		//总拉起订单金额
		mysql.DB.Table("prepaid_phone_orders").Select("sum(account_orders) as all_pull_order_amount").Scan(&Data)
		//总成功订单金额
		mysql.DB.Table("prepaid_phone_orders").Where("status =?", 2).Select("sum(account_practical) as all_pull_order_amount_and_success").Scan(&Data)
		if Data.AllPullOrderCount == 0 {
			Data.AllSuccessPer = 0
		} else {
			Data.AllSuccessPer = float64(Data.AllPullOrderCountAndSuccess / Data.AllPullOrderCount)
		}

		tools.ReturnError200Data(c, Data, "success")
		return
	}

}
