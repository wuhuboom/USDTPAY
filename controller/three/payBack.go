package three

import (
	"encoding/base64"
	"encoding/json"
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"io/ioutil"
	"strconv"
	"strings"
)

// GetPayInformationBack 支付成功回调
func GetPayInformationBack(c *gin.Context) {

	body, _ := ioutil.ReadAll(c.Request.Body)
	if body != nil {
		fmt.Println(string(body))
	}

	var jsonDataTwo ReturnBase64
	err := c.BindJSON(&jsonDataTwo)
	if err != nil {
		tools.ReturnError101(c, err.Error())
		return
	}
	var apiKey = viper.GetString("project.ApiKey")
	if tools.ApiSign([]byte(apiKey), []byte(jsonDataTwo.Data), []byte(apiKey)) != jsonDataTwo.Sign {
		tools.ReturnError101(c, "非法请求")
		return
	}
	sDec, err1 := base64.StdEncoding.DecodeString(jsonDataTwo.Data)
	if err1 != nil {
		tools.ReturnError101(c, "非法请求")
		return
	}
	fmt.Println(string(sDec))
	zap.L().Debug("GetPayInformationBack:" + string(sDec))
	var jsonData GetPayInformationBackData
	err = json.Unmarshal(sDec, &jsonData)
	if err != nil {
		tools.ReturnError101(c, "非法请求")
		return
	}

	//保存日志
	BackContent := map[string]string{}
	BackContent["data"] = jsonDataTwo.Data
	BackContent["sign"] = jsonDataTwo.Sign
	BackContentJson, _ := json.Marshal(BackContent)
	log := model.BackLog{Kinds: 1, BackContent: string(BackContentJson), JsonContent: string(sDec), TxHash: jsonData.Data.TxHash}
	defer log.CreatedBackLog(mysql.DB)
	//  余额清零
	if jsonData.Type == "balance" {
		var jsonDataTwo BalanceType
		err = json.Unmarshal(sDec, &jsonDataTwo)
		if err != nil {
			tools.ReturnError101(c, "非法请求")
			return
		}
		zap.L().Debug("余额变动,用户:" + jsonDataTwo.Data.Addr)
		re := model.ReceiveAddress{Username: jsonDataTwo.Data.User}
		re.UpdateReceiveAddressLastInformationTo0(mysql.DB)
		log.Kinds = 2
		tools.ReturnError200(c, "余额变动成功")
		return
	}

	//判断这个  TxHash是否已经被使用过?
	orders := model.PrepaidPhoneOrders{TxHash: jsonData.Data.TxHash}
	if orders.IfUseThisTxHash(mysql.DB) {
		tools.ReturnError200(c, "TxHash已被使用")
		return
	}

	rare := model.ReceiveAddress{}
	rare.Kinds = 1
	mysql.DB.Where("address=?", jsonData.Data.To).First(&rare)
	p1 := model.PrepaidPhoneOrders{
		Username:          jsonData.Data.UserID,
		Successfully:      jsonData.Data.Timestamp / 1000,
		RechargeType:      strings.ToUpper(jsonData.Data.Token),
		RechargeAddress:   jsonData.Data.To, //收账地址
		CollectionAddress: jsonData.Data.From,
		TxHash:            jsonData.Data.TxHash,
	} //玩家地址

	acc := strconv.Itoa(jsonData.Data.Amount)
	p1.AccountPractical, _ = tools.ToDecimal(acc, 6).Float64()
	if rare.Kinds == 1 {
		//寻找这个账号最早地充值订单
		p1.UpdateMaxCreatedOfStatusToTwo(mysql.DB, viper.GetInt64("project.OrderEffectivityTime"))
	} else {
		//池子的地址
		p1.UpdatePondOrderCratedAndUpdated(mysql.DB)
	}

	//更新钱包地址
	newMoney, _ := tools.ToDecimal(jsonData.Data.Balance, 6).Float64()
	R := model.ReceiveAddress{LastGetAccount: p1.AccountPractical, Username: jsonData.Data.UserID, Updated: jsonData.Data.Timestamp / 1000, Money: newMoney}
	R.UpdateReceiveAddressLastInformation(mysql.DB)
	//更新总的账变
	change := model.BalanceChange{OriginalAmount: 0, ChangeAmount: p1.AccountPractical, NowAmount: 0}
	change.Add(mysql.DB)
	tools.ReturnError200(c, "插入成功")
	return
}
