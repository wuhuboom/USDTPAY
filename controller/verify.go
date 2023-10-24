package controller

// LoginVerify 登录参数
type LoginVerify struct {
	Username     string `form:"username"  binding:"required"`
	Password     string `form:"password"  binding:"required"`
	GoogleCode   string `form:"googleCode"  binding:"omitempty,max=6" `
	GoogleSecret string `form:"googleSecret"  binding:"omitempty" `
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

type ConsoleManagementData struct {
	TodayPullOrderCount            int64   `json:"today_pull_order_count"`              //今日拉单总数
	TodayPullOrderCountAndSuccess  int64   `json:"today_pull_order_count_and_success"`  //今日成功支付笔数
	TodayPullOrderAmount           float64 `json:"today_pull_order_amount"`             //今日拉单总金额
	TodayPullOrderAmountAndSuccess float64 `json:"today_pull_order_amount_and_success"` //今日收取金额
	TodaySuccessPer                float64 `json:"today_success_per"`                   //今日订单支付成功率
	AllPullOrderCount              int64   `json:"all_pull_order_count"`                //总拉单总数
	AllPullOrderCountAndSuccess    int64   `json:"all_pull_order_count_and_success"`    //总成功支付笔数
	AllPullOrderAmount             float64 `json:"all_pull_order_amount"`               //总拉单总金额
	AllPullOrderAmountAndSuccess   float64 `json:"all_pull_order_amount_and_success"`   //总收取金额
	AllSuccessPer                  float64 `json:"all_success_per"`                     //总成功率
}
