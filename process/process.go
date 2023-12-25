package process

import (
	"context"
	"encoding/json"
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/btcsuite/btcutil/base58"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
	"gorm.io/gorm"
	"io"
	"io/ioutil"
	"math"
	"math/big"
	"net/http"
	"strings"
	"time"
)

func CratedPoolAddress(db *gorm.DB) {
	for true {

		if viper.GetString("Project.ModelA") == "no" {
			return
		}
		var pondSize int64
		db.Model(&model.ReceiveAddress{}).Where("kinds=?", 2).Count(&pondSize) //3
		//admin := model.Admin{}
		config := model.Config{}
		config.GetConfig(db)
		db.Where("id=?", 1).First(&config)
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
			req.Header.Set("TRON-PRO-API-KEY", viper.GetString("project.TronApiKey"))
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

func Repair(db *gorm.DB) {
	pp := make([]model.PrepaidPhoneOrders, 0)
	db.Find(&pp)
	for _, i2 := range pp {
		t := time.Unix(i2.Created, 0)
		db.Model(&model.PrepaidPhoneOrders{}).Where("id=?", i2.ID).Updates(&model.PrepaidPhoneOrders{
			Date: t.Format("2006-01-02"),
		})
		fmt.Println("--")
	}

	fmt.Println("修复完毕")
}

// BlockCheck 区块订单检查 到账检查
func BlockCheck(db *gorm.DB, block int) {
	Kn := 1
	too := 0
	start := 0
	for i := 0; i < Kn; i++ {
		url := fmt.Sprintf("https://apilist.tronscanapi.com/api/transaction?"+
			"sort=-timestamp&count=true&limit=50&start=%d&block=%d", start, block)
		req, err := http.NewRequest("GET", url, nil)
		req.Header.Set("TRON-PRO-API-KEY", viper.GetString("project.TronApiKey"))
		if err != nil {
			fmt.Println(err.Error(), 223)
			continue
		}
		req.Header.Add("accept", "application/json")
		res, err := http.DefaultClient.Do(req)
		if err != nil {
			fmt.Println(err.Error(), 229)
			continue
		}
		body, err := ioutil.ReadAll(res.Body)
		if err != nil {
			fmt.Println(err.Error(), 234)
			continue
		}
		var tt1 TransactionData
		err = json.Unmarshal(body, &tt1)
		if err != nil {
			fmt.Println(err.Error(), 240)
			continue
		}
		fmt.Println(tt1.Total)
		for _, datum := range tt1.Data {
			if datum.TriggerInfo.ContractAddress == "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" {
				address := model.ReceiveAddress{Address: datum.TriggerInfo.Parameter.To}
				if bb, _ := address.IsExist(db); bb == true {
					//判断hash是否被使用过
					//判断这个  TxHash是否已经被使用过?
					orders := model.PrepaidPhoneOrders{TxHash: datum.Hash}
					if orders.IfUseThisTxHash(db) == false {
						p1 := model.PrepaidPhoneOrders{
							Username:          address.Username,
							Successfully:      datum.Timestamp / 1000,
							RechargeType:      "USDT",
							RechargeAddress:   datum.TriggerInfo.Parameter.To, //收账地址
							CollectionAddress: datum.OwnerAddress,
							TxHash:            datum.Hash,
						} //玩家地址
						acc := datum.TriggerInfo.Parameter.Value
						p1.AccountPractical, _ = tools.ToDecimal(acc, 6).Float64()
						if address.Kinds == 1 {
							//寻找这个账号最早地充值订单
							p1.UpdateMaxCreatedOfStatusToTwo(mysql.DB, viper.GetInt64("project.OrderEffectivityTime"))
						} else {
							//池子的地址
							p1.UpdatePondOrderCratedAndUpdated(mysql.DB)
						}
						//更新钱包地址
						newMoney, _ := tools.ToDecimal(datum.TriggerInfo.Parameter.Value, 6).Float64()
						R := model.ReceiveAddress{
							LastGetAccount: p1.AccountPractical,
							Username:       address.Username,
							Updated:        datum.Timestamp / 1000,
							Money:          newMoney}
						R.UpdateReceiveAddressLastInformation(db)
						//更新总的账变
						change := model.BalanceChange{OriginalAmount: 0, ChangeAmount: p1.AccountPractical, NowAmount: 0}
						change.Add(db)
					}
				}
				fmt.Println(datum.TriggerInfo.Parameter.To)
				fmt.Println(tools.ToDecimal(datum.TriggerInfo.Parameter.Value, 6).Float64())
				too++
			}
		}
		start = start + 50
		if start < tt1.Total {
			Kn++
		}

		time.Sleep(time.Second * 1)

	}

}

type TransactionData struct {
	Total      int `json:"total"`
	RangeTotal int `json:"rangeTotal"`
	Data       []struct {
		Block         int      `json:"block"`
		Hash          string   `json:"hash"`
		Timestamp     int64    `json:"timestamp"`
		OwnerAddress  string   `json:"ownerAddress"`
		ToAddressList []string `json:"toAddressList"`
		ToAddress     string   `json:"toAddress"`
		ContractType  int      `json:"contractType"`
		Confirmed     bool     `json:"confirmed"`
		Revert        bool     `json:"revert"`
		ContractData  struct {
			Data            string `json:"data,omitempty"`
			OwnerAddress    string `json:"owner_address"`
			ContractAddress string `json:"contract_address,omitempty"`
			Amount          int64  `json:"amount,omitempty"`
			ToAddress       string `json:"to_address,omitempty"`
			AssetName       string `json:"asset_name,omitempty"`
			TokenInfo       struct {
				TokenId      string `json:"tokenId"`
				TokenAbbr    string `json:"tokenAbbr"`
				TokenName    string `json:"tokenName"`
				TokenDecimal int    `json:"tokenDecimal"`
				TokenCanShow int    `json:"tokenCanShow"`
				TokenType    string `json:"tokenType"`
				TokenLogo    string `json:"tokenLogo"`
				TokenLevel   string `json:"tokenLevel"`
				Vip          bool   `json:"vip"`
			} `json:"tokenInfo,omitempty"`
			AccountAddress  string `json:"account_address,omitempty"`
			Balance         int64  `json:"balance,omitempty"`
			Resource        string `json:"resource,omitempty"`
			ReceiverAddress string `json:"receiver_address,omitempty"`
		} `json:"contractData"`
		SmartCalls  string `json:"SmartCalls"`
		Events      string `json:"Events"`
		Id          string `json:"id"`
		Data        string `json:"data"`
		Fee         string `json:"fee"`
		ContractRet string `json:"contractRet"`
		Result      string `json:"result"`
		Amount      string `json:"amount"`
		CheatStatus bool   `json:"cheatStatus"`
		Cost        struct {
			NetFee             int `json:"net_fee"`
			EnergyPenaltyTotal int `json:"energy_penalty_total"`
			EnergyUsage        int `json:"energy_usage"`
			Fee                int `json:"fee"`
			EnergyFee          int `json:"energy_fee"`
			EnergyUsageTotal   int `json:"energy_usage_total"`
			OriginEnergyUsage  int `json:"origin_energy_usage"`
			NetUsage           int `json:"net_usage"`
		} `json:"cost"`
		TokenInfo struct {
			TokenId      string `json:"tokenId"`
			TokenAbbr    string `json:"tokenAbbr"`
			TokenName    string `json:"tokenName"`
			TokenDecimal int    `json:"tokenDecimal"`
			TokenCanShow int    `json:"tokenCanShow"`
			TokenType    string `json:"tokenType"`
			TokenLogo    string `json:"tokenLogo"`
			TokenLevel   string `json:"tokenLevel"`
			Vip          bool   `json:"vip"`
		} `json:"tokenInfo"`
		TokenType   string `json:"tokenType"`
		TriggerInfo struct {
			Method    string `json:"method"`
			Data      string `json:"data"`
			Parameter struct {
				Value string `json:"_value"`
				To    string `json:"_to"`
				From  string `json:"_from,omitempty"`
			} `json:"parameter"`
			MethodId        string `json:"methodId"`
			MethodName      string `json:"methodName"`
			ContractAddress string `json:"contract_address"`
			CallValue       int    `json:"call_value"`
		} `json:"trigger_info,omitempty"`
		RiskTransaction bool `json:"riskTransaction"`
	} `json:"data"`
	WholeChainTxCount int64 `json:"wholeChainTxCount"`
	ContractMap       struct {
		TYBpSnC8MDBP4WNLxFWP7GqkLvB3PF27Sq bool `json:"TYBpSnC8mDBP4wNLxFWP7GqkLvB3pF27Sq"`
		TJkRacqx52Zr6TTqhsJvBM9ECnLMUz2QSe bool `json:"TJkRacqx52zr6TTqhsJvBM9eCnLMUz2qSe"`
		THLSe8YVzXrCEMyX6HCPrGfDW3USskejUJ bool `json:"THLSe8YVzXrCEMyX6hCPrGfDW3USskejUJ"`
		TS1P19KeAW4GRHgqT6HNoUYiXZicoaLKgi bool `json:"TS1p19keAW4GRHgqT6hNoUYiXZicoaLKgi"`
		TLzcufNPP6JABLLMpxp8CfdDPL3DLrounH bool `json:"TLzcufNPP6jABLLMpxp8CfdDPL3dLrounH"`
		TXj3Zjdd7KAWLUHKZCD3XpJSpVqWxPGPtA bool `json:"TXj3Zjdd7kAWLUHKZCD3XpJSpVqWxPGPtA"`
		TDm9PLda9AjBQGu45R5RdvgM4JKXWvWHZf bool `json:"TDm9pLda9AjBQGu45r5RdvgM4JKXWvWHZf"`
		TXgG1Sti1LrmRaMJ7EZG1HDjmGno9YMNQ2 bool `json:"TXgG1Sti1LrmRaMJ7eZG1hDjmGno9YMNQ2"`
		TJkbXK6XQ4E8IKczJ9ZRZ5QQFmaR9RdDJq bool `json:"TJkbXK6XQ4E8iKczJ9ZRZ5qQFmaR9RdDJq"`
		TLVYmSEFdkt4HehUqmheAMzuRZiFdnsKFE bool `json:"TLVYmSEFdkt4HehUqmheAMzuRZiFdnsKFE"`
		TTdf3D8KPM4U6JkExdLGfNqHZ4JoeK9Exs bool `json:"TTdf3d8KPM4u6JkExdLGfNqHZ4JoeK9Exs"`
		TJwpVgCjcQWirxhvD3NMcr2UmpVZRKsNEc bool `json:"TJwpVgCjcQWirxhvD3NMcr2umpVZRKsNEc"`
		TVxHqdFazwgKJmKoMGZwbnVuKjPXt9D7Wa bool `json:"TVxHqdFazwgKJmKoMGZwbnVuKjPXt9D7wa"`
		TE4KVBSj6FFvUZG1VwYuRzcNtJhXxWxXNL bool `json:"TE4KVBSj6fFvUZG1VwYuRzcNtJhXxWxXNL"`
		TMHhK3TeBPbPZDQQeJDjEZwhU6Tfy4GUyK bool `json:"TMHhK3TeBPbPZDQQeJDjEZwhU6Tfy4gUyK"`
		TW2KrdvmrstL1Vgtp5WGzeq1Ru4XPH86Xg bool `json:"TW2krdvmrstL1vgtp5wGzeq1Ru4XPH86Xg"`
		TJVDxwgPdp7GGoDnHs9GfGphHuGFjKLMF8 bool `json:"TJVDxwgPdp7GGoDnHs9GfGphHuGFjKLMF8"`
		TBXHCjR6HT1BT1Qf3DGq7WaDnXkvuQAXoU bool `json:"TBXHCjR6hT1BT1qf3dGq7WaDnXkvuQAXoU"`
		TPKBoRBp4UjU72GNX87Nkar1GKn1UjdP93 bool `json:"TPKBoRBp4ujU72gNX87Nkar1GKn1ujdP93"`
		TXkXkmTfPgrSN7NnM3MeY7YbNwHr36KsBz bool `json:"TXkXkmTfPgrSN7nnM3MeY7YbNwHr36ksBz"`
		TRLwmN7WD5WWf3X2CpnBgkEpKWQbT1TLji bool `json:"TRLwmN7wD5wWf3x2CpnBgkEpKWQbT1tLji"`
		TES1PgfoTE8FMRBWeyJgxqmRB2Uss74NVm bool `json:"TES1PgfoTE8fMRBWeyJgxqmRB2uss74nVm"`
		TU8Y5MTWq7No7UC1Ca7KBdYfk3XVj81GXx bool `json:"TU8y5MTWq7no7UC1ca7KBdYfk3XVj81gXx"`
		TTpfN9TBnGJJ7THUt7YaWv1MdtGG3FGhDB bool `json:"TTpfN9TBnGJJ7tHUt7YaWv1MdtGG3FGhDB"`
		TLPiduJ775ImccG8Z3C2OLHhtBxoAaDiHV bool `json:"TLPiduJ775imccG8z3c2oLHhtBxoAaDiHV"`
		TWNUDrHzGnkwWLic7ZAztiW1Bm6DaC9VAi bool `json:"TWNUDrHzGnkwWLic7ZAztiW1bm6daC9vAi"`
		TGvXaV3Tdy5Ze3B9DK8CE3FnfQG9XLN8Hq bool `json:"TGvXaV3Tdy5Ze3b9DK8CE3fnfQG9xLN8Hq"`
		TWbpLsrjtg1CxrYRJ3LSbg1NhWPvK8UG83 bool `json:"TWbpLsrjtg1cxrYRJ3LSbg1NhWPvK8uG83"`
		TVtLDTC94T8QCn6JA3ZpPrUyYyFGCepQJF bool `json:"TVtLDTC94T8QCn6JA3zpPrUyYyFGCepQJF"`
		TWvZBMQPGwQmkRh5Dx54Ry8AfLuKqhX5Od bool `json:"TWvZBMQPGwQmkRh5dx54Ry8AfLuKqhX5od"`
		TVxJZRKRaTLwzV5LSa5WD536BXzNyC52J6 bool `json:"TVxJZRKRaTLwzV5LSa5WD536BXzNyC52j6"`
		TYQoLXmPuN3DVWGWsQqKfNfRSRBgU3EGd1 bool `json:"TYQoLXmPuN3dVWGWsQqKfNfRSRBgU3eGd1"`
		TYerzAjpSezJd4HchdYbestRvXmzK24ZLW bool `json:"TYerzAjpSezJd4HchdYbestRvXmzK24zLW"`
		TWqjtfvcK6JxSLetYL7QJHWrxSRLS55J1J bool `json:"TWqjtfvcK6jxSLetYL7qJHWrxSRLS55J1J"`
		TY8TJAZxxNWD6XSgt9JkgfbbbJ6F5RQ8CS bool `json:"TY8tJAZxxNWD6xSgt9JkgfbbbJ6f5rQ8CS"`
		TEcopCBxx2ObRvNwkjRaonGj3AjxA31HFS bool `json:"TEcopCBxx2obRvNwkjRaonGj3AjxA31HFS"`
		TDys4JWbF2Q1RRS3FnZ6S7TZCfx7BriXMh bool `json:"TDys4jWbF2Q1RRS3fnZ6s7TZCfx7briXMh"`
		TX4U3DEQ5MZPSan9GZzBG7ZUbMbnHnMvkQ bool `json:"TX4U3DEQ5mZPSan9gZzBG7ZUbMbnHnMvkQ"`
		TTFNwrBML7IdvwR36GkCQP9TWwkzmegoA4 bool `json:"TTFNwrBML7idvwR36GkCQP9tWwkzmegoA4"`
		TWVbxU33Sd4IK7M6Ryej6U6Mr1SwMqtVD5 bool `json:"TWVbxU33Sd4iK7m6Ryej6u6mr1SwMqtVD5"`
		TLMhWGSS7C33Tg15JmCFWNS37BnT1RQ15K bool `json:"TLMhWGSS7c33Tg15jmCFWNS37BnT1RQ15K"`
		TTBwUETaAuoK8UsSeUsT4N7ZZQotgokEDX bool `json:"TTBwUETaAuoK8UsSeUsT4n7zZQotgokEDX"`
		TUuSB4ETir1TEvjpurcemhtBsYDKigUPAb bool `json:"TUuSB4eTir1tEvjpurcemhtBsYDKigUPAb"`
		TEn8HKKUMSSPnbP9IUxFZUt4NWVyf7XLhU bool `json:"TEn8HKKUMSSPnbP9iUxFZUt4NWVyf7xLhU"`
		TTDktDK9Ribh2G1SwwdzywbYbLSC25VaAQ bool `json:"TTDktDK9ribh2g1SwwdzywbYbLSC25VaAQ"`
		TEi2HVWDRMo61PAoi1Dwbn8HNXufkwEVyp bool `json:"TEi2hVWDRMo61PAoi1Dwbn8hNXufkwEVyp"`
		TY9X2FdTJgoJYBfDZQ8REMeRAsWtchoWuN bool `json:"TY9x2FdTJgoJYBfDZQ8rEMeRAsWtchoWuN"`
		TA4CGA5Sb9PSRV76SKFZKYBiv66Xh8Gztg bool `json:"TA4cGA5sb9PSRV76sKFZKYBiv66xh8gztg"`
		TU1DkdF3N9GNHHFnoQDV6F9Wn8TCE7Dhwp bool `json:"TU1DkdF3n9GNHHFnoQDV6F9Wn8TCE7dhwp"`
		TS8D9SDGc81HoFHMYStG7FVGNumx2A1REa bool `json:"TS8D9sDGc81HoFHMYStG7FVGNumx2a1REa"`
		TWrTTVizHkF3LxTiWE9VCnhWC7CBpNXe22 bool `json:"TWrTTVizHkF3LxTiWE9vCnhWC7CBpNXe22"`
		TVFbe8XYaD7KY4KeMasEH2JvW683Kfc4Jx bool `json:"TVFbe8xYaD7KY4keMasEH2JvW683kfc4Jx"`
		TXUYrpjqxZCcAPygUZTC8APeZZyjwGoUoo bool `json:"TXUYrpjqxZCcAPygUZTC8aPeZZyjwGoUoo"`
		TEN4KrL95T6CSWZwb71GaiXj5ZbadJuT3O bool `json:"TEN4KrL95t6cSWZwb71gaiXj5ZbadJuT3o"`
		TUh8OSRV6FQzdPTmLJiVWgSL5Hy25KK86H bool `json:"TUh8oSRV6fQzdPTmLJiVWgSL5hy25kK86H"`
		TDmuz2WDWyPTbwGAmBojz5LcTxPrRvgVTP bool `json:"TDmuz2WDWyPTbwGAmBojz5LcTxPrRvgVTP"`
		TRXEPcW2Fe5F8CtbRGjqGoKiZ2YLppVmNw bool `json:"TRXEPcW2fe5f8CtbRGjqGoKiZ2yLppVmNw"`
		TEdyLGA7QQ7I8VCbb4DbEbAe9ByV4Y7ZBx bool `json:"TEdyLGA7qQ7i8VCbb4DbEbAe9byV4y7ZBx"`
		TXEFhXTidasbSTxWKmCZ8HciPd86G643ZL bool `json:"TXEFhXTidasbSTxWKmCZ8HciPd86g643ZL"`
		TQ4Xs49CKpBXWBnfKsF36LoALRnjF4Grwv bool `json:"TQ4xs49CKpBXWBnfKsF36LoALRnjF4Grwv"`
		TBD1Fk7QEBMeEA2Qh48ItFsHSkVgYPBKMv bool `json:"TBD1fk7qEBMeEA2qh48itFsHSkVgYPBKMv"`
		TR7NHqjeKQxGTCi8Q8ZY4PL8OtSzgjLj6T bool `json:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"`
		TLwnLa9HoZbT9SZ6ECsHvGVF67L78UoGZ5 bool `json:"TLwnLa9hoZbT9sZ6eCsHvGVF67L78UoGZ5"`
		TAdjY1Z1RNGiC4GeXsqGeDyANP8IA24Ht7 bool `json:"TAdjY1Z1RNGiC4geXsqGeDyANP8iA24ht7"`
		TM4Y3Z2J8A4Q3W64XoUpfQTYqvR8GDx1HY bool `json:"TM4Y3Z2J8a4q3w64XoUpfQTYqvR8GDx1hY"`
		TMcvYzWHLtcKmLu528Xcrxrf25YMJ9ZeaK bool `json:"TMcvYzWHLtcKmLu528xcrxrf25yMJ9zeaK"`
		TE9DjVnxeLyKSf1276FskUe1NQ997Eqxgd bool `json:"TE9djVnxeLyKSf1276fskUe1nQ997Eqxgd"`
		TPjpJwDe4KZ2THHAL5AUEPLvVoH9TN1Eq4 bool `json:"TPjpJwDe4kZ2THHAL5AUEPLvVoH9TN1eq4"`
		TAjTGgQsrYzyi3AnUhQkLen76ZTfU88888 bool `json:"TAjTGgQsrYzyi3anUhQkLen76ZTfU88888"`
		TMQm6AD4CpDEWDceMuAweYmDcYihB95GC3 bool `json:"TMQm6AD4cpDEWDceMuAweYmDcYihB95gC3"`
		TW5JCNKqf6G4BuoAuXEjzHj3CP44ZnSSNa bool `json:"TW5JCNKqf6g4BuoAuXEjzHj3cP44znSSNa"`
		TDCRKSTX4S6AdzhcQDnPkamBDbGvPLcssE bool `json:"TDCRKSTX4s6AdzhcQDnPkamBDbGvPLcssE"`
		TRrAhAkLvAF7RiMQoJMLjHq51U9Q1OhU7D bool `json:"TRrAhAkLvAF7RiMQoJMLjHq51u9Q1ohU7D"`
	} `json:"contractMap"`
	ContractInfo struct {
		TR7NHqjeKQxGTCi8Q8ZY4PL8OtSzgjLj6T struct {
			IsToken bool   `json:"isToken"`
			Tag1    string `json:"tag1"`
			Tag1Url string `json:"tag1Url"`
			Name    string `json:"name"`
			Risk    bool   `json:"risk"`
			Vip     bool   `json:"vip"`
		} `json:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"`
	} `json:"contractInfo"`
	NormalAddressInfo struct {
		TYBpSnC8MDBP4WNLxFWP7GqkLvB3PF27Sq struct {
			Risk bool `json:"risk"`
		} `json:"TYBpSnC8mDBP4wNLxFWP7GqkLvB3pF27Sq"`
		TJkRacqx52Zr6TTqhsJvBM9ECnLMUz2QSe struct {
			Risk bool `json:"risk"`
		} `json:"TJkRacqx52zr6TTqhsJvBM9eCnLMUz2qSe"`
		THLSe8YVzXrCEMyX6HCPrGfDW3USskejUJ struct {
			Risk bool `json:"risk"`
		} `json:"THLSe8YVzXrCEMyX6hCPrGfDW3USskejUJ"`
		TS1P19KeAW4GRHgqT6HNoUYiXZicoaLKgi struct {
			Risk bool `json:"risk"`
		} `json:"TS1p19keAW4GRHgqT6hNoUYiXZicoaLKgi"`
		TLzcufNPP6JABLLMpxp8CfdDPL3DLrounH struct {
			Risk bool `json:"risk"`
		} `json:"TLzcufNPP6jABLLMpxp8CfdDPL3dLrounH"`
		TXj3Zjdd7KAWLUHKZCD3XpJSpVqWxPGPtA struct {
			Risk bool `json:"risk"`
		} `json:"TXj3Zjdd7kAWLUHKZCD3XpJSpVqWxPGPtA"`
		TDm9PLda9AjBQGu45R5RdvgM4JKXWvWHZf struct {
			Risk bool `json:"risk"`
		} `json:"TDm9pLda9AjBQGu45r5RdvgM4JKXWvWHZf"`
		TXgG1Sti1LrmRaMJ7EZG1HDjmGno9YMNQ2 struct {
			Risk bool `json:"risk"`
		} `json:"TXgG1Sti1LrmRaMJ7eZG1hDjmGno9YMNQ2"`
		TJkbXK6XQ4E8IKczJ9ZRZ5QQFmaR9RdDJq struct {
			Risk bool `json:"risk"`
		} `json:"TJkbXK6XQ4E8iKczJ9ZRZ5qQFmaR9RdDJq"`
		TLVYmSEFdkt4HehUqmheAMzuRZiFdnsKFE struct {
			Risk bool `json:"risk"`
		} `json:"TLVYmSEFdkt4HehUqmheAMzuRZiFdnsKFE"`
		TTdf3D8KPM4U6JkExdLGfNqHZ4JoeK9Exs struct {
			Risk bool `json:"risk"`
		} `json:"TTdf3d8KPM4u6JkExdLGfNqHZ4JoeK9Exs"`
		TJwpVgCjcQWirxhvD3NMcr2UmpVZRKsNEc struct {
			Risk bool `json:"risk"`
		} `json:"TJwpVgCjcQWirxhvD3NMcr2umpVZRKsNEc"`
		TVxHqdFazwgKJmKoMGZwbnVuKjPXt9D7Wa struct {
			Risk bool `json:"risk"`
		} `json:"TVxHqdFazwgKJmKoMGZwbnVuKjPXt9D7wa"`
		TE4KVBSj6FFvUZG1VwYuRzcNtJhXxWxXNL struct {
			Risk bool `json:"risk"`
		} `json:"TE4KVBSj6fFvUZG1VwYuRzcNtJhXxWxXNL"`
		TMHhK3TeBPbPZDQQeJDjEZwhU6Tfy4GUyK struct {
			Risk bool `json:"risk"`
		} `json:"TMHhK3TeBPbPZDQQeJDjEZwhU6Tfy4gUyK"`
		TW2KrdvmrstL1Vgtp5WGzeq1Ru4XPH86Xg struct {
			Risk bool `json:"risk"`
		} `json:"TW2krdvmrstL1vgtp5wGzeq1Ru4XPH86Xg"`
		TJVDxwgPdp7GGoDnHs9GfGphHuGFjKLMF8 struct {
			Risk bool `json:"risk"`
		} `json:"TJVDxwgPdp7GGoDnHs9GfGphHuGFjKLMF8"`
		TBXHCjR6HT1BT1Qf3DGq7WaDnXkvuQAXoU struct {
			Risk bool `json:"risk"`
		} `json:"TBXHCjR6hT1BT1qf3dGq7WaDnXkvuQAXoU"`
		TPKBoRBp4UjU72GNX87Nkar1GKn1UjdP93 struct {
			Risk bool `json:"risk"`
		} `json:"TPKBoRBp4ujU72gNX87Nkar1GKn1ujdP93"`
		TXkXkmTfPgrSN7NnM3MeY7YbNwHr36KsBz struct {
			Risk bool `json:"risk"`
		} `json:"TXkXkmTfPgrSN7nnM3MeY7YbNwHr36ksBz"`
		TRLwmN7WD5WWf3X2CpnBgkEpKWQbT1TLji struct {
			Risk bool `json:"risk"`
		} `json:"TRLwmN7wD5wWf3x2CpnBgkEpKWQbT1tLji"`
		TES1PgfoTE8FMRBWeyJgxqmRB2Uss74NVm struct {
			Risk bool `json:"risk"`
		} `json:"TES1PgfoTE8fMRBWeyJgxqmRB2uss74nVm"`
		TU8Y5MTWq7No7UC1Ca7KBdYfk3XVj81GXx struct {
			Risk bool `json:"risk"`
		} `json:"TU8y5MTWq7no7UC1ca7KBdYfk3XVj81gXx"`
		TTpfN9TBnGJJ7THUt7YaWv1MdtGG3FGhDB struct {
			Risk bool `json:"risk"`
		} `json:"TTpfN9TBnGJJ7tHUt7YaWv1MdtGG3FGhDB"`
		TLPiduJ775ImccG8Z3C2OLHhtBxoAaDiHV struct {
			Risk bool `json:"risk"`
		} `json:"TLPiduJ775imccG8z3c2oLHhtBxoAaDiHV"`
		TWNUDrHzGnkwWLic7ZAztiW1Bm6DaC9VAi struct {
			Risk bool `json:"risk"`
		} `json:"TWNUDrHzGnkwWLic7ZAztiW1bm6daC9vAi"`
		TGvXaV3Tdy5Ze3B9DK8CE3FnfQG9XLN8Hq struct {
			Risk bool `json:"risk"`
		} `json:"TGvXaV3Tdy5Ze3b9DK8CE3fnfQG9xLN8Hq"`
		TWbpLsrjtg1CxrYRJ3LSbg1NhWPvK8UG83 struct {
			Risk bool `json:"risk"`
		} `json:"TWbpLsrjtg1cxrYRJ3LSbg1NhWPvK8uG83"`
		TVtLDTC94T8QCn6JA3ZpPrUyYyFGCepQJF struct {
			Risk bool `json:"risk"`
		} `json:"TVtLDTC94T8QCn6JA3zpPrUyYyFGCepQJF"`
		TWvZBMQPGwQmkRh5Dx54Ry8AfLuKqhX5Od struct {
			Risk bool `json:"risk"`
		} `json:"TWvZBMQPGwQmkRh5dx54Ry8AfLuKqhX5od"`
		TVxJZRKRaTLwzV5LSa5WD536BXzNyC52J6 struct {
			Risk bool `json:"risk"`
		} `json:"TVxJZRKRaTLwzV5LSa5WD536BXzNyC52j6"`
		TYQoLXmPuN3DVWGWsQqKfNfRSRBgU3EGd1 struct {
			Risk bool `json:"risk"`
		} `json:"TYQoLXmPuN3dVWGWsQqKfNfRSRBgU3eGd1"`
		TYerzAjpSezJd4HchdYbestRvXmzK24ZLW struct {
			Risk bool `json:"risk"`
		} `json:"TYerzAjpSezJd4HchdYbestRvXmzK24zLW"`
		TWqjtfvcK6JxSLetYL7QJHWrxSRLS55J1J struct {
			Risk bool `json:"risk"`
		} `json:"TWqjtfvcK6jxSLetYL7qJHWrxSRLS55J1J"`
		TY8TJAZxxNWD6XSgt9JkgfbbbJ6F5RQ8CS struct {
			Risk bool `json:"risk"`
		} `json:"TY8tJAZxxNWD6xSgt9JkgfbbbJ6f5rQ8CS"`
		TEcopCBxx2ObRvNwkjRaonGj3AjxA31HFS struct {
			Risk bool `json:"risk"`
		} `json:"TEcopCBxx2obRvNwkjRaonGj3AjxA31HFS"`
		TDys4JWbF2Q1RRS3FnZ6S7TZCfx7BriXMh struct {
			Risk bool `json:"risk"`
		} `json:"TDys4jWbF2Q1RRS3fnZ6s7TZCfx7briXMh"`
		TX4U3DEQ5MZPSan9GZzBG7ZUbMbnHnMvkQ struct {
			Risk bool `json:"risk"`
		} `json:"TX4U3DEQ5mZPSan9gZzBG7ZUbMbnHnMvkQ"`
		TTFNwrBML7IdvwR36GkCQP9TWwkzmegoA4 struct {
			Risk bool `json:"risk"`
		} `json:"TTFNwrBML7idvwR36GkCQP9tWwkzmegoA4"`
		TWVbxU33Sd4IK7M6Ryej6U6Mr1SwMqtVD5 struct {
			Risk bool `json:"risk"`
		} `json:"TWVbxU33Sd4iK7m6Ryej6u6mr1SwMqtVD5"`
		TLMhWGSS7C33Tg15JmCFWNS37BnT1RQ15K struct {
			Risk bool `json:"risk"`
		} `json:"TLMhWGSS7c33Tg15jmCFWNS37BnT1RQ15K"`
		TTBwUETaAuoK8UsSeUsT4N7ZZQotgokEDX struct {
			Risk bool `json:"risk"`
		} `json:"TTBwUETaAuoK8UsSeUsT4n7zZQotgokEDX"`
		TUuSB4ETir1TEvjpurcemhtBsYDKigUPAb struct {
			Risk bool `json:"risk"`
		} `json:"TUuSB4eTir1tEvjpurcemhtBsYDKigUPAb"`
		TEn8HKKUMSSPnbP9IUxFZUt4NWVyf7XLhU struct {
			Risk bool `json:"risk"`
		} `json:"TEn8HKKUMSSPnbP9iUxFZUt4NWVyf7xLhU"`
		TTDktDK9Ribh2G1SwwdzywbYbLSC25VaAQ struct {
			Risk bool `json:"risk"`
		} `json:"TTDktDK9ribh2g1SwwdzywbYbLSC25VaAQ"`
		TEi2HVWDRMo61PAoi1Dwbn8HNXufkwEVyp struct {
			Risk bool `json:"risk"`
		} `json:"TEi2hVWDRMo61PAoi1Dwbn8hNXufkwEVyp"`
		TY9X2FdTJgoJYBfDZQ8REMeRAsWtchoWuN struct {
			Risk bool `json:"risk"`
		} `json:"TY9x2FdTJgoJYBfDZQ8rEMeRAsWtchoWuN"`
		TA4CGA5Sb9PSRV76SKFZKYBiv66Xh8Gztg struct {
			Risk bool `json:"risk"`
		} `json:"TA4cGA5sb9PSRV76sKFZKYBiv66xh8gztg"`
		TU1DkdF3N9GNHHFnoQDV6F9Wn8TCE7Dhwp struct {
			Risk bool `json:"risk"`
		} `json:"TU1DkdF3n9GNHHFnoQDV6F9Wn8TCE7dhwp"`
		TS8D9SDGc81HoFHMYStG7FVGNumx2A1REa struct {
			Risk bool `json:"risk"`
		} `json:"TS8D9sDGc81HoFHMYStG7FVGNumx2a1REa"`
		TWrTTVizHkF3LxTiWE9VCnhWC7CBpNXe22 struct {
			Risk bool `json:"risk"`
		} `json:"TWrTTVizHkF3LxTiWE9vCnhWC7CBpNXe22"`
		TVFbe8XYaD7KY4KeMasEH2JvW683Kfc4Jx struct {
			Risk bool `json:"risk"`
		} `json:"TVFbe8xYaD7KY4keMasEH2JvW683kfc4Jx"`
		TXUYrpjqxZCcAPygUZTC8APeZZyjwGoUoo struct {
			Risk bool `json:"risk"`
		} `json:"TXUYrpjqxZCcAPygUZTC8aPeZZyjwGoUoo"`
		TEN4KrL95T6CSWZwb71GaiXj5ZbadJuT3O struct {
			Risk bool `json:"risk"`
		} `json:"TEN4KrL95t6cSWZwb71gaiXj5ZbadJuT3o"`
		TUh8OSRV6FQzdPTmLJiVWgSL5Hy25KK86H struct {
			Risk bool `json:"risk"`
		} `json:"TUh8oSRV6fQzdPTmLJiVWgSL5hy25kK86H"`
		TDmuz2WDWyPTbwGAmBojz5LcTxPrRvgVTP struct {
			Risk bool `json:"risk"`
		} `json:"TDmuz2WDWyPTbwGAmBojz5LcTxPrRvgVTP"`
		TRXEPcW2Fe5F8CtbRGjqGoKiZ2YLppVmNw struct {
			Risk bool `json:"risk"`
		} `json:"TRXEPcW2fe5f8CtbRGjqGoKiZ2yLppVmNw"`
		TEdyLGA7QQ7I8VCbb4DbEbAe9ByV4Y7ZBx struct {
			Risk bool `json:"risk"`
		} `json:"TEdyLGA7qQ7i8VCbb4DbEbAe9byV4y7ZBx"`
		TXEFhXTidasbSTxWKmCZ8HciPd86G643ZL struct {
			Risk bool `json:"risk"`
		} `json:"TXEFhXTidasbSTxWKmCZ8HciPd86g643ZL"`
		TQ4Xs49CKpBXWBnfKsF36LoALRnjF4Grwv struct {
			Risk bool `json:"risk"`
		} `json:"TQ4xs49CKpBXWBnfKsF36LoALRnjF4Grwv"`
		TBD1Fk7QEBMeEA2Qh48ItFsHSkVgYPBKMv struct {
			Risk bool `json:"risk"`
		} `json:"TBD1fk7qEBMeEA2qh48itFsHSkVgYPBKMv"`
		TR7NHqjeKQxGTCi8Q8ZY4PL8OtSzgjLj6T struct {
			Risk bool `json:"risk"`
		} `json:"TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t"`
		TLwnLa9HoZbT9SZ6ECsHvGVF67L78UoGZ5 struct {
			Risk bool `json:"risk"`
		} `json:"TLwnLa9hoZbT9sZ6eCsHvGVF67L78UoGZ5"`
		TAdjY1Z1RNGiC4GeXsqGeDyANP8IA24Ht7 struct {
			Risk bool `json:"risk"`
		} `json:"TAdjY1Z1RNGiC4geXsqGeDyANP8iA24ht7"`
		TM4Y3Z2J8A4Q3W64XoUpfQTYqvR8GDx1HY struct {
			Risk bool `json:"risk"`
		} `json:"TM4Y3Z2J8a4q3w64XoUpfQTYqvR8GDx1hY"`
		TMcvYzWHLtcKmLu528Xcrxrf25YMJ9ZeaK struct {
			Risk bool `json:"risk"`
		} `json:"TMcvYzWHLtcKmLu528xcrxrf25yMJ9zeaK"`
		TE9DjVnxeLyKSf1276FskUe1NQ997Eqxgd struct {
			Risk bool `json:"risk"`
		} `json:"TE9djVnxeLyKSf1276fskUe1nQ997Eqxgd"`
		TPjpJwDe4KZ2THHAL5AUEPLvVoH9TN1Eq4 struct {
			Risk bool `json:"risk"`
		} `json:"TPjpJwDe4kZ2THHAL5AUEPLvVoH9TN1eq4"`
		TAjTGgQsrYzyi3AnUhQkLen76ZTfU88888 struct {
			Risk bool `json:"risk"`
		} `json:"TAjTGgQsrYzyi3anUhQkLen76ZTfU88888"`
		TMQm6AD4CpDEWDceMuAweYmDcYihB95GC3 struct {
			Risk bool `json:"risk"`
		} `json:"TMQm6AD4cpDEWDceMuAweYmDcYihB95gC3"`
		TW5JCNKqf6G4BuoAuXEjzHj3CP44ZnSSNa struct {
			Risk bool `json:"risk"`
		} `json:"TW5JCNKqf6g4BuoAuXEjzHj3cP44znSSNa"`
		TDCRKSTX4S6AdzhcQDnPkamBDbGvPLcssE struct {
			Risk bool `json:"risk"`
		} `json:"TDCRKSTX4s6AdzhcQDnPkamBDbGvPLcssE"`
		TRrAhAkLvAF7RiMQoJMLjHq51U9Q1OhU7D struct {
			Risk bool `json:"risk"`
		} `json:"TRrAhAkLvAF7RiMQoJMLjHq51u9Q1ohU7D"`
	} `json:"normalAddressInfo"`
}

func BlockCheckFromTo(db *gorm.DB, client *redis.Client) {
	var ctx = context.Background()
	for true {
		//获取区块
		Block, _ := client.Get(ctx, "Block").Int()
		fmt.Println(fmt.Sprintf("正在检查区块:%d", Block))
		BlockCheck(db, Block)
		if Block < GetNewBlock() {
			c := Block + 1
			client.Set(ctx, "Block", c, 0)
		}
	}
}

// GetNewBlock 获取最新区块
func GetNewBlock() int {
	url := fmt.Sprintf("https://apilist.tronscanapi.com/api/block?limit=1")
	req, err := http.NewRequest("GET", url, nil)
	req.Header.Set("TRON-PRO-API-KEY", viper.GetString("project.TronApiKey"))
	if err != nil {
		fmt.Println(err.Error(), 708)
		return 0
	}
	req.Header.Add("accept", "application/json")
	res, err := http.DefaultClient.Do(req)
	if err != nil {
		fmt.Println(err.Error(), 714)
		return 0
	}
	body, err := ioutil.ReadAll(res.Body)
	if err != nil {
		fmt.Println(err.Error(), 719)
		return 0
	}
	var tt1 GetNewBlockData
	err = json.Unmarshal(body, &tt1)
	if err != nil {
		fmt.Println(err.Error(), 725)
		return 0
	}
	return tt1.Data[0].Number
}

type GetNewBlockData struct {
	Total      int `json:"total"`
	RangeTotal int `json:"rangeTotal"`
	Data       []struct {
		Number         int     `json:"number"`
		Hash           string  `json:"hash"`
		Size           int     `json:"size"`
		Timestamp      int64   `json:"timestamp"`
		TxTrieRoot     string  `json:"txTrieRoot"`
		ParentHash     string  `json:"parentHash"`
		WitnessId      int     `json:"witnessId"`
		WitnessAddress string  `json:"witnessAddress"`
		NrOfTrx        int     `json:"nrOfTrx"`
		WitnessName    string  `json:"witnessName"`
		Version        string  `json:"version"`
		Fee            float64 `json:"fee"`
		Confirmed      bool    `json:"confirmed"`
		Confirmations  int     `json:"confirmations"`
		NetUsage       int     `json:"netUsage"`
		EnergyUsage    int     `json:"energyUsage"`
		BlockReward    int     `json:"blockReward"`
		VoteReward     int     `json:"voteReward"`
		Revert         bool    `json:"revert"`
	} `json:"data"`
}

var UpdateBalanceProcessChan = make(chan int)

func UpdateBalanceProcess(redis *redis.Client) {
	c := context.Background()
	for {
		add, err := redis.RPop(c, "updateMoney").Result()
		if err == nil {
			var re model.ReceiveAddress
			err := json.Unmarshal([]byte(add), &re)
			if err != nil {
				break
			}
			//UpdateBalance(re)

			UpdateBalance2(re)
			//fmt.Println(re)
		}
		time.Sleep(time.Second * 1)

	}
}

func UpdateBalance(v model.ReceiveAddress) {
	url := "https://apilist.tronscanapi.com/api/account/tokens?address=" + v.Address + "&start=0&limit=20&token=&hidden=0&show=0&sortType=0"
	req, _ := http.NewRequest("GET", url, nil)
	req.Header.Add("accept", "application/json")
	req.Header.Set("TRON-PRO-API-KEY", viper.GetString("project.TronApiKey"))
	res, _ := http.DefaultClient.Do(req)
	body, _ := ioutil.ReadAll(res.Body)
	fmt.Println(string(body))
	var tt2 Ta2
	err := json.Unmarshal(body, &tt2)
	if err != nil {
		return
	}
	var newMoney float64
	newMoney = 0
	for _, datum := range tt2.Data {
		if datum.TokenId == "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t" {
			newMoney, _ = tools.ToDecimal(datum.Balance, 6).Float64()
		}
	}
	//fmt.Printf("余额:%f", tt1.Data[0].AssetInUsd)
	//usd := ToDecimal(arrayA[1], 6)
	////更新数据
	ups := make(map[string]interface{})
	ups["Money"] = newMoney
	ups["Updated"] = time.Now().Unix()
	err = mysql.DB.Model(model.ReceiveAddress{}).Where("id=?", v.ID).Updates(ups).Error
	fmt.Println(newMoney)
	//调动 余额变动
	if math.Abs(newMoney-v.Money) > 1 {
		change := model.AccountChange{ChangeAmount: math.Abs(newMoney - v.Money), Kinds: 1, OriginalAmount: v.Money, NowAmount: newMoney, ReceiveAddressName: v.Username}
		change.Add(mysql.DB)
	}
	if err != nil {
		fmt.Println("更新失败")
	}
	time.Sleep(1 * time.Second)
}

func UpdateBalance2(v model.ReceiveAddress) {
	url := "https://api.trongrid.io/wallet/triggersmartcontract"
	hexAddress, err := base58ToHex(v.Address)
	if err != nil {
		fmt.Println("Error:", err)
		return
	}
	jsonData := `{
	   "visible": true,
	   "call_value": 0,
	   "function_selector": "balanceOf(address)",
	   "owner_address": "` + v.Address + `",
	   "fee_limit": 1000000000,
	   "contract_address": "TR7NHqjeKQxGTCi8q8ZY4pL8otSzgjLj6t",
	   "parameter": "00000000` + hexAddress + `"
	}`
	//fmt.Println(jsonData)
	payload := strings.NewReader(jsonData)
	req, _ := http.NewRequest("POST", url, payload)
	req.Header.Add("accept", "application/json")
	req.Header.Add("content-type", "application/json")
	req.Header.Add("TRON-PRO-API-KEY", "6efbfebb-c0c4-4363-ab2b-96b7e50b694a")
	res, _ := http.DefaultClient.Do(req)
	defer res.Body.Close()
	body, _ := io.ReadAll(res.Body)

	var ud UpdateBalance2Data
	err = json.Unmarshal(body, &ud)
	if err != nil {
		fmt.Println("852" + err.Error())
		return
	}
	fmt.Println(v.Address)
	fmt.Println(ud.ConstantResult[0])
	decimal, err := toDecimal(ud.ConstantResult[0])
	if err != nil {
		fmt.Println("Failed to convert hex to decimal:", err)
		return
	}
	divisor := big.NewInt(1000000)
	//fmt.Println(new(big.Int).Div(decimal, divisor))
	newMoney, _ := new(big.Int).Div(decimal, divisor).Float64()
	ups := make(map[string]interface{})
	ups["Money"] = newMoney
	ups["Updated"] = time.Now().Unix()
	err = mysql.DB.Model(model.ReceiveAddress{}).Where("id=?", v.ID).Updates(ups).Error
	fmt.Println(newMoney)
	//调动 余额变动
	if math.Abs(newMoney-v.Money) > 1 {
		change := model.AccountChange{ChangeAmount: math.Abs(newMoney - v.Money), Kinds: 1, OriginalAmount: v.Money, NowAmount: newMoney, ReceiveAddressName: v.Username}
		change.Add(mysql.DB)
	}
	if err != nil {
		fmt.Println("更新失败")
	}

}

type UpdateBalance2Data struct {
	Result struct {
		Result bool `json:"result"`
	} `json:"result"`
	EnergyUsed     int      `json:"energy_used"`
	ConstantResult []string `json:"constant_result"`
	EnergyPenalty  int      `json:"energy_penalty"`
	Transaction    struct {
		Ret []struct {
		} `json:"ret"`
		Visible bool   `json:"visible"`
		TxID    string `json:"txID"`
		RawData struct {
			Contract []struct {
				Parameter struct {
					Value struct {
						Data            string `json:"data"`
						OwnerAddress    string `json:"owner_address"`
						ContractAddress string `json:"contract_address"`
					} `json:"value"`
					TypeUrl string `json:"type_url"`
				} `json:"parameter"`
				Type string `json:"type"`
			} `json:"contract"`
			RefBlockBytes string `json:"ref_block_bytes"`
			RefBlockHash  string `json:"ref_block_hash"`
			Expiration    int64  `json:"expiration"`
			FeeLimit      int    `json:"fee_limit"`
			Timestamp     int64  `json:"timestamp"`
		} `json:"raw_data"`
		RawDataHex string `json:"raw_data_hex"`
	} `json:"transaction"`
}

func base58ToHex(base58Addr string) (string, error) {
	//00000000000000411d44658e2acbaf0b1f0c2ccb69904f470b6b27a16dbd8abe
	decoded := base58.Decode(base58Addr)
	hexAddr := fmt.Sprintf("%064x", decoded)
	return hexAddr[0 : len(hexAddr)-8], nil
}
func toDecimal(hexString string) (*big.Int, error) {
	hexString = strings.TrimPrefix(hexString, "0x")
	decimal := new(big.Int)
	_, success := decimal.SetString(hexString, 16)
	if !success {
		return nil, fmt.Errorf("Failed to convert hex string to decimal")
	}
	return decimal, nil
}
