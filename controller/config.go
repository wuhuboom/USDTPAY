package controller

import (
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
	"strconv"
)

// Config 配置
func Config(c *gin.Context) {
	action := c.PostForm("action")
	//查看
	if action == "check" {
		config := model.Config{}
		mysql.DB.Where("id=?", 1).First(&config)
		tools.ReturnError200Data(c, config, "OK")
		return
	}
	//修改
	if action == "update" {
		maxPond, _ := strconv.Atoi(c.PostForm("max_pond"))
		expiration, _ := strconv.ParseInt(c.PostForm("expiration"), 10, 64)
		pondAmount, _ := strconv.ParseFloat(c.PostForm("pond_amount"), 64)
		if maxPond == 0 {
			maxPond = 1000
		}
		if expiration == 0 {
			expiration = 30
		}
		if pondAmount == 0 {
			pondAmount = 5
		}
		//池设置要单独处理下
		var total int64
		mysql.DB.Model(&model.ReceiveAddress{}).Where("kinds=? and status=?", 2, 1).Count(&total)
		// 现在所有的地址大于池地址
		if total > int64(maxPond) {
			race := make([]model.ReceiveAddress, 0)
			mysql.DB.Where("kinds=? and status=?", 2, 1).
				Order("receive_nums asc ").Limit(int(total) - maxPond).Find(&race)
			for _, address := range race {
				mysql.DB.Model(&model.ReceiveAddress{}).Where("id=?", address.ID).Updates(&model.ReceiveAddress{Status: 2})
			}
		} else {
			race := make([]model.ReceiveAddress, 0)
			mysql.DB.Where("kinds=? and status=?", 2, 2).
				Order("receive_nums asc ").Limit(int(total) - maxPond).Find(&race)
			for _, address := range race {
				mysql.DB.Model(&model.ReceiveAddress{}).Where("id=?", address.ID).Updates(&model.ReceiveAddress{Status: 1})
			}
		}

		mysql.DB.Model(&model.Config{}).Where("id=?", 1).Updates(&model.Config{
			MaxPond:    maxPond,
			Expiration: expiration,
			PondAmount: pondAmount,
		})
		tools.ReturnError200(c, "OK")
		return
	}

}
