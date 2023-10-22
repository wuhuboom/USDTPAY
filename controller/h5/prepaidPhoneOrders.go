package h5

import (
	"encoding/base64"
	"encoding/json"
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"fmt"
	"github.com/spf13/viper"
	"math/rand"
	"strconv"
	"strings"
	"time"

	"example.com/m/tools"
	"example.com/m/util"
	"github.com/gin-gonic/gin"
	"go.uber.org/zap"
)

// CreatePrepaidPhoneOrders 生成充值订单
func CreatePrepaidPhoneOrders(c *gin.Context) {
	var jsonDataT CreatePrepaidPhoneOrdersT
	err := c.BindJSON(&jsonDataT)
	if err != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 14, err.Error()))
		tools.ReturnError101(c, "err1:"+err.Error())
		return
	}
	decodeString, err1 := base64.StdEncoding.DecodeString(jsonDataT.Data)
	if err1 != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 22, err1.Error()))
		tools.ReturnError101(c, "err2:"+err1.Error())
		return
	}
	origData, err2 := util.RsaDecrypt(decodeString)
	if err2 != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 30, err2.Error()))
		tools.ReturnError101(c, "err3:"+err2.Error())
		return
	}
	var jsonData CreatePrepaidPhoneOrdersData
	err3 := json.Unmarshal(origData, &jsonData)
	if err3 != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 38, err3.Error()))
		tools.ReturnError101(c, "err4:"+err3.Error())
		return
	}
	//判断用户是存在充值地址
	if jsonData.Username == "" {
		tools.ReturnError101(c, "用户名不可以为空")
		return
	}
	//判断充值类型
	if strings.ToUpper(jsonData.RechargeType) != "USDT" {
		tools.ReturnError101(c, "RechargeType is error")
		return
	}
	//判断平台订单是否重复
	affected := mysql.DB.Where("platform_order=?", jsonData.PlatformOrder).Limit(1).Find(&model.PrepaidPhoneOrders{}).RowsAffected
	if affected == 1 {
		tools.ReturnError101(c, "不要重复提交")
		return
	}

	//判断这个玩家是否有专属地址
	re := model.ReceiveAddress{}
	rowsAffected := mysql.DB.Where("username=?", jsonData.Username).Limit(1).Find(&re).RowsAffected
	if rowsAffected == 0 {
		//用户不存在  判断金额 是去池子里面获取 还是生成专属地址
		config := model.Config{}
		config.GetConfig(mysql.DB)
		//获取配置文件
		fmt.Println(config)
		// 订单金额小于设置金额   要从池中获取地址
		if jsonData.AccountOrders <= config.PondAmount {
			//判断是否有已在使用的地址并且没有过期
			pp := model.PrepaidPhoneOrders{}
			err := mysql.DB.
				Where("username=? and created >  ? and account_orders <= ?",
					jsonData.Username, time.Now().Unix()-config.Expiration*60, config.PondAmount).
				First(&pp).Error

			if err == nil {
				//地址还未过期重复使用这个地址
				mysql.DB.Where("address=?", pp.RechargeAddress).First(&re)
				//更新地址库这个地址的最后一次使用时间
				mysql.DB.Model(&model.ReceiveAddress{}).Where("address=?", pp.RechargeAddress).Updates(
					&model.ReceiveAddress{
						LastUseTime: time.Now().Unix() + config.Expiration*60,
						ReceiveNums: re.ReceiveNums + 1})
			} else {
				//随机获取一个最新的地址给玩家
				GetAddressLock.Lock()
				err = mysql.DB.
					Where("kinds=? and last_use_time < ? and status=?",
						2, time.Now().Unix(), 1).
					Order("receive_nums asc").First(&re).Error

				//池里面没有地址
				if err != nil {
					//寻找没有开启的地址
					err = mysql.DB.Where("kinds=? and last_use_time < ? and status=?",
						2, time.Now().Unix(), 2).
						Order("receive_nums asc").First(&re).Error
					if err != nil {
						//如果池地址是满的状态要提现管理员添加地址  后期做成飞机警报
						//model.Log{Ips: c.ClientIP(),Content: ""}
						tools.ReturnError101(c, "There are no more addresses in the pool")
						GetAddressLock.Unlock()
						return
					}
				} else {
					//更新地址的最新
					mysql.DB.Model(&model.ReceiveAddress{}).Where("id=?", re.ID).Updates(
						&model.ReceiveAddress{
							LastUseTime: time.Now().Unix() + config.Expiration*60,
							ReceiveNums: re.ReceiveNums + 1})
					GetAddressLock.Unlock()
					//生成地址日志
					logs := model.Log{Kinds: 3, Ips: c.ClientIP(), Content: fmt.Sprintf("玩家:%s,成功获取地址:%s", jsonData.Username, re.Address)}
					logs.CreatedLogs(mysql.DB)
				}

			}
		} else {
			//大于这个金额生成 专属地址
			re.Username = jsonData.Username
			re.CreateUsername(mysql.DB, viper.GetString("project.ThreeUrl"))
			if re.Address == "" {
				tools.ReturnError101(c, "返回空的地址,稍后重试-1")
				//生成 地址日志
				log := model.Log{Ips: c.ClientIP(),
					Content: fmt.Sprintf("用户:%s,获取专属地址失败,获取连接:%s,ip:%s", re.Username, viper.GetString("eth.ThreeUrl"), c.ClientIP()),
					Kinds:   2}
				log.CreatedLogs(mysql.DB)
				return
			}
			//生成地址日志
			logs := model.Log{Kinds: 3, Ips: c.ClientIP(), Content: fmt.Sprintf("玩家:%s,成功获取地址:%s", jsonData.Username, re.Address)}
			logs.CreatedLogs(mysql.DB)
		}
	}
	//生成充值订单
	p := model.PrepaidPhoneOrders{
		PlatformOrder:   jsonData.PlatformOrder,
		RechargeAddress: re.Address,
		AccountOrders:   jsonData.AccountOrders,
		Username:        jsonData.Username,
		RechargeType:    jsonData.RechargeType,
		BackUrl:         jsonData.BackUrl,
		ThreeOrder:      time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000)),
		Status:          1, Created: time.Now().Unix()}
	err = mysql.DB.Save(&p).Error
	if err != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 88, err.Error()))
		tools.ReturnError101(c, "err6:"+err.Error())
		return
	}
	aUrl := viper.GetString("project.RechargingJumpAddress") + "?Address=" + re.Address + "&RechargeType=" + jsonData.RechargeType + "&AccountOrders=" + strconv.FormatFloat(jsonData.AccountOrders, 'f', 2, 64) + "&PlatformOrder=" + jsonData.PlatformOrder
	var oo ReturnData
	oo.UrlAddress = aUrl
	data, err := json.Marshal(oo)
	if err != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 99, err.Error()))
		tools.ReturnError101(c, "err7:"+err.Error())
		return
	}
	data, err = util.RsaEncryptForEveryOne(data)
	if err != nil {
		zap.L().Debug(tools.SetPrintStyle("CreatePrepaidPhoneOrders", 105, err.Error()))
		tools.ReturnError101(c, "err8:"+err.Error())
		return
	}
	//充值订单创建成功
	tools.ReturnError200Data(c, base64.StdEncoding.EncodeToString(data), "Ok")
	return
}
