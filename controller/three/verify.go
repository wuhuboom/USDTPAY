package three

type ReturnBase64 struct {
	Data string `json:"data"`
	Sign string `json:"sign"`
}

type Data struct {
	TxHash      string `json:"txHash" binding:"required"`
	BlockNumber int    `json:"blockNumber" binding:"required"`
	Timestamp   int64  `json:"timestamp" binding:"required"`
	From        string `json:"from" binding:"required"`
	To          string `json:"to" binding:"required"`
	Amount      int    `json:"amount" binding:"required"`
	Token       string `json:"token" binding:"required"`
	UserID      string `json:"userId" binding:"required"`
	Balance     string `json:"balance" binding:"required"`
}

type GetPayInformationBackData struct {
	Type string `json:"type" binding:"required"`
	Data Data   `json:"data" binding:"required"`
	Sign string `json:"sign" binding:"required"`
}

type BalanceType struct {
	Data struct {
		Addr    string `json:"addr"`
		Balance string `json:"balance"`
		Seq     int64  `json:"seq"`
		User    string `json:"user"`
	} `json:"data"`
	Type string `json:"type"`
}
