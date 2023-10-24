package process

import (
	"encoding/json"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io/ioutil"
	"net/http"
	"time"
)

func CratedPoolAddress(db *gorm.DB) {
	for true {
		var pondSize int64
		db.Model(&model.ReceiveAddress{}).Where("kinds=?", 2).Count(&pondSize) //3
		admin := model.Admin{}
		config := model.Config{}
		config.GetConfig(db)
		db.Where("id=?", 1).First(&admin)
		if int(pondSize) < config.MaxPond {
			//池的地址  < 配置地址数量
			needPool := config.MaxPond - int(pondSize)
			for i := 0; i < needPool; i++ {
				re := model.ReceiveAddress{}
				re.Username = "PoolConnection" + time.Now().Format("2006-01-02") + string(tools.RandString(8))
				re.Kinds = 2
				re.ReceiveNums = 0
				fmt.Println(viper.GetString("project.ThreeUrl"))
				re.CreateUsername(db, viper.GetString("project.ThreeUrl"))
				if re.Address == "" {
					//生成 地址日志
					log := model.Log{Ips: "127.0.0.1",
						Content: fmt.Sprintf("用户:%s,获取专属地址失败,获取连接:%s,ip:%s", re.Username, viper.GetString("eth.ThreeUrl"), "127.0.0.1"),
						Kinds:   2}
					log.CreatedLogs(db)
				}
			}
		}
		time.Sleep(time.Minute)
	}
}

func Test(db *gorm.DB) {

	for i := 0; i < 1000; i++ {

		re := model.ReceiveAddress{}
		err := db.
			Where("kinds=? and last_use_time < ?",
				2, time.Now().Unix()).
			Order("receive_nums asc").First(&re).Error

		if err == nil {
			db.Model(&model.ReceiveAddress{}).Where("id=?", re.ID).Updates(&model.ReceiveAddress{
				ReceiveNums: re.ReceiveNums + 1,
				LastUseTime: time.Now().Unix() + 60*3600,
			})
		}
		fmt.Println(i)
	}

}

type Ta struct {
	Total        int `json:"total"`
	ContractInfo struct {
	} `json:"contractInfo"`
	RangeTotal     int `json:"rangeTotal"`
	TokenTransfers []struct {
		TransactionId  string `json:"transaction_id"`
		BlockTs        int64  `json:"block_ts"`
		FromAddress    string `json:"from_address"`
		FromAddressTag struct {
			FromAddressTag     string `json:"from_address_tag"`
			FromAddressTagLogo string `json:"from_address_tag_logo"`
		} `json:"from_address_tag"`
		ToAddress    string `json:"to_address"`
		ToAddressTag struct {
			ToAddressTagLogo string `json:"to_address_tag_logo"`
			ToAddressTag     string `json:"to_address_tag"`
		} `json:"to_address_tag"`
		Block           int    `json:"block"`
		ContractAddress string `json:"contract_address"`
		TriggerInfo     struct {
			Method    string `json:"method"`
			Data      string `json:"data"`
			Parameter struct {
				Value string `json:"_value"`
				To    string `json:"_to"`
			} `json:"parameter"`
			MethodName      string `json:"methodName"`
			ContractAddress string `json:"contract_address"`
			CallValue       int    `json:"call_value"`
		} `json:"trigger_info"`
		Quant          string `json:"quant"`
		ApprovalAmount string `json:"approval_amount"`
		EventType      string `json:"event_type"`
		ContractType   string `json:"contract_type"`
		Confirmed      bool   `json:"confirmed"`
		ContractRet    string `json:"contractRet"`
		FinalResult    string `json:"finalResult"`
		TokenInfo      struct {
			TokenId      string `json:"tokenId"`
			TokenAbbr    string `json:"tokenAbbr"`
			TokenName    string `json:"tokenName"`
			TokenDecimal int    `json:"tokenDecimal"`
			TokenCanShow int    `json:"tokenCanShow"`
			TokenType    string `json:"tokenType"`
			TokenLogo    string `json:"tokenLogo"`
			TokenLevel   string `json:"tokenLevel"`
			IssuerAddr   string `json:"issuerAddr"`
			Vip          bool   `json:"vip"`
		} `json:"tokenInfo"`
		FromAddressIsContract bool `json:"fromAddressIsContract"`
		ToAddressIsContract   bool `json:"toAddressIsContract"`
		Revert                bool `json:"revert"`
	} `json:"token_transfers"`
}

type Ta2 struct {
	Total int `json:"total"`
	Data  []struct {
		Amount           interface{} `json:"amount"`
		Quantity         interface{} `json:"quantity"`
		TokenId          string      `json:"tokenId"`
		TokenPriceInUsd  float64     `json:"tokenPriceInUsd"`
		TokenName        string      `json:"tokenName"`
		TokenAbbr        string      `json:"tokenAbbr"`
		TokenCanShow     int         `json:"tokenCanShow"`
		TokenLogo        string      `json:"tokenLogo"`
		TokenPriceInTrx  float64     `json:"tokenPriceInTrx"`
		AmountInUsd      float64     `json:"amountInUsd"`
		Balance          string      `json:"balance"`
		TokenDecimal     int         `json:"tokenDecimal"`
		TokenType        string      `json:"tokenType"`
		Vip              bool        `json:"vip"`
		NrOfTokenHolders int         `json:"nrOfTokenHolders,omitempty"`
		TransferCount    int         `json:"transferCount,omitempty"`
		Project          string      `json:"project,omitempty"`
	} `json:"data"`
	ContractMap struct {
		TR7NHqjeKQxGTCi8Q8ZY4PL8OtSzgjLj6T bool `json:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"`
		Field2                             bool `json:"_"`
	} `json:"contractMap"`
	ContractInfo struct {
		TR7NHqjeKQxGTCi8Q8ZY4PL8OtSzgjLj6T struct {
			Tag1    string `json:"tag1"`
			Tag1Url string `json:"tag1Url"`
			Name    string `json:"name"`
			Vip     bool   `json:"vip"`
		} `json:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"`
	} `json:"contractInfo"`
}

func CheckLastGetMoneyTime(db *gorm.DB) {
	for true {
		rA := make([]model.ReceiveAddress, 0)
		db.Find(&rA)
		for _, address := range rA {
			url := "https://apilist.tronscanapi.com/api/token_trc20/transfers?limit=20&start=0&sort=-timestamp&count=true&relatedAddress=" + address.Address
			req, err := http.NewRequest("GET", url, nil)
			req.Header.Set("TRON-PRO-API-KEY", viper.GetString("app.TronApiKey"))
			if err != nil {
				continue
			}
			req.Header.Add("accept", "application/json")
			res, err := http.DefaultClient.Do(req)
			if err != nil {
				continue
			}
			body, err := ioutil.ReadAll(res.Body)
			if err != nil {
				continue
			}
			var tt1 Ta
			err = json.Unmarshal(body, &tt1)
			if err != nil {
				continue
			}
			if len(tt1.TokenTransfers) > 0 {
				//最后一次接收转账的时间
				db.Model(&model.ReceiveAddress{}).Where("id=?", address.ID).Updates(&model.ReceiveAddress{TheLastGetMoneyTime: tt1.TokenTransfers[0].BlockTs})
			}

			//	fmt.Println("检查地址: " + address.Address + "完毕")
		}

		time.Sleep(time.Minute * 60 * 24)
	}
}
