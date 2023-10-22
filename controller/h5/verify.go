package h5

import "sync"

type CreatePrepaidPhoneOrdersT struct {
	Data string `json:"data"`
}

type CreatePrepaidPhoneOrdersData struct {
	PlatformOrder string  `json:"PlatformOrder" binding:"required"`  //平台订单号
	Username      string  `json:"Username" binding:"required"`       //充值用户名
	AccountOrders float64 `json:"AccountOrders"  binding:"required"` //充值金额
	RechargeType  string  `json:"RechargeType"  binding:"required"`  //充值类型
	BackUrl       string  `json:"BackUrl"  binding:"required" `      // 回调地址
}

type ReturnData struct {
	UrlAddress string
}

var GetAddressLock sync.RWMutex
