package controller

import (
	"encoding/base64"
	"encoding/json"
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"example.com/m/util"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"sync"
)

var orderBackLock sync.RWMutex

// TopUpOrder 充值订单管理
func TopUpOrder(c *gin.Context) {
	action := c.PostForm("action")
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.PrepaidPhoneOrders, 0)
		Db := mysql.DB
		var total int64

		type ACC struct {
			AccountOrders    float64 `json:"account_orders"`
			AccountPractical float64 `json:"account_practical"`
		}

		var bcc ACC
		// 用户名
		if content, isExist := c.GetPostForm("username"); isExist == true {
			Db = Db.Where("username=?", content)
		}

		//平台订单号
		if content, isExist := c.GetPostForm("platform_order"); isExist == true {
			Db = Db.Where("platform_order=?", content)
		}

		//三方平台订单号
		if content, isExist := c.GetPostForm("three_order"); isExist == true {
			Db = Db.Where("three_order=?", content)
		}

		//收账地址
		if content, isExist := c.GetPostForm("recharge_address"); isExist == true {
			Db = Db.Where("recharge_address=?", content)
		}

		//订单状态
		if content, isExist := c.GetPostForm("status"); isExist == true {
			Db = Db.Where("status=?", content)
		}

		//是否回调
		if content, isExist := c.GetPostForm("three_back"); isExist == true {
			Db = Db.Where("three_back=?", content)
		}

		//日期条件
		if start, isExist := c.GetPostForm("start_time"); isExist == true {
			if end, isExist := c.GetPostForm("end_time"); isExist == true {
				Db = Db.Where("successfully >= ?", start).Where("successfully<=?", end)
				mysql.DB.Raw("select sum(account_orders) as account_orders ,sum(account_practical) as  account_practical from prepaid_phone_orders where successfully  BETWEEN ? AND ?", start, end).Scan(&bcc)

			}
		}
		Db.Table("prepaid_phone_orders").Count(&total)
		Db = Db.Model(&model.PrepaidPhoneOrders{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}

		c.JSON(http.StatusOK, gin.H{
			"code":             0,
			"count":            total,
			"data":             role,
			"accountOrders":    bcc.AccountOrders,
			"accountPractical": bcc.AccountPractical,
		})
		return
	}
	if action == "update" {
		remark := c.PostForm("remark")
		id := c.PostForm("id")
		if remark == "" || id == "" {
			tools.ReturnError101(c, "remark is not  null")
			return
		}
		mysql.DB.Model(&model.PrepaidPhoneOrders{}).Where("id=?", id).Updates(&model.PrepaidPhoneOrders{Remark: remark})
		tools.ReturnError200(c, "OK")
	}
	//订单回调
	if action == "orderBack" {
		orderBackLock.Lock()
		defer orderBackLock.Unlock()
		id := c.PostForm("id")
		txHash := c.PostForm("tx_hash")
		if len(txHash) != 64 {
			tools.ReturnError101(c, "填写正确的hash值")
			return
		}
		actualAmount, _ := strconv.ParseFloat(c.PostForm("actual_amount"), 64)
		p := model.PrepaidPhoneOrders{}
		affected := mysql.DB.Where("id=?", id).Limit(1).Find(&p).RowsAffected
		if affected == 0 {
			tools.ReturnError101(c, "订单不存在")
			return
		}
		if p.ThreeBack != 1 {
			tools.ReturnError101(c, "不要重复回调")
			return
		}
		//判断这个hash 是否存在
		p2 := model.PrepaidPhoneOrders{}
		if p2.IfUseThisTxHash(mysql.DB) {
			tools.ReturnError101(c, "这个hash已经使用过")
			return
		}
		rowsAffected := mysql.DB.Model(&model.PrepaidPhoneOrders{}).Where("id=?", id).Updates(&model.PrepaidPhoneOrders{
			ThreeBack:        4,
			Status:           2,
			AccountPractical: actualAmount}).RowsAffected
		if rowsAffected == 0 {
			tools.ReturnError101(c, "更新失败")
			return
		}

		//回调给三方
		type Create struct {
			PlatformOrder    string
			RechargeAddress  string
			Username         string
			AccountOrders    float64 //订单充值金额
			AccountPractical float64 //  实际充值的金额
			RechargeType     string
			BackUrl          string
		}
		var tt Create
		tt.PlatformOrder = p.PlatformOrder
		tt.RechargeAddress = p.RechargeAddress
		tt.Username = p.Username
		tt.AccountOrders = p.AccountOrders
		tt.AccountPractical = actualAmount
		tt.RechargeType = p.RechargeType
		data, err := json.Marshal(tt)
		if err != nil {
			tools.ReturnError101(c, err.Error())
			return
		}
		data, err = util.RsaEncryptForEveryOne(data)
		util.BackUrlToPay(p.BackUrl, base64.StdEncoding.EncodeToString(data))
		tools.ReturnError200(c, "回调成功")
		return
	}
}
