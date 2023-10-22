package tools

import (
	"encoding/json"
	"fmt"
	"github.com/agclqq/goencryption"
	"io/ioutil"
	"net/http"
	"testing"
	"time"
)

func GET() {
	priv, err := goencryption.GenPrvKey(2048)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", priv)
	pub, err := goencryption.GenPubKeyFromPrvKey(priv)
	if err != nil {
		fmt.Println(err)
	}
	fmt.Printf("%s\n", pub)
}

func TestApiSign(t *testing.T) {

	url := "https://api.shasta.trongrid.io/wallet/getaccount"

	url = "https://apilist.tronscanapi.com/api/account/token_asset_overview?address=TSs1bE2PaNahbMqi9yctTZ6wZ7d2Fbzm8N"

	req, _ := http.NewRequest("GET", url, nil)

	req.Header.Add("accept", "application/json")

	res, _ := http.DefaultClient.Do(req)

	defer res.Body.Close()
	body, _ := ioutil.ReadAll(res.Body)

	//fmt.Println(res)
	fmt.Println(string(body))

	var tt1 Ta
	err := json.Unmarshal(body, &tt1)
	if err != nil {
		return
	}

	fmt.Println(tt1.Data[0].AssetInUsd)

}

type Ta struct {
	TotalAssetInTrx float64 `json:"totalAssetInTrx"`
	Data            []struct {
		TokenId         string  `json:"tokenId"`
		TokenName       string  `json:"tokenName"`
		TokenAbbr       string  `json:"tokenAbbr"`
		TokenDecimal    int     `json:"tokenDecimal"`
		TokenCanShow    int     `json:"tokenCanShow"`
		TokenType       string  `json:"tokenType"`
		TokenLogo       string  `json:"tokenLogo"`
		Vip             bool    `json:"vip"`
		Balance         string  `json:"balance"`
		TokenPriceInTrx float64 `json:"tokenPriceInTrx"`
		TokenPriceInUsd float64 `json:"tokenPriceInUsd"`
		AssetInTrx      float64 `json:"assetInTrx"`
		AssetInUsd      float64 `json:"assetInUsd"`
		Percent         float64 `json:"percent"`
	} `json:"data"`
	TotalAssetInUsd float64 `json:"totalAssetInUsd"`
}

func TestName(t *testing.T) {
	req := make(map[string]interface{})
	req["user"] = "wangyi"
	req["ts"] = time.Now().UnixMilli()

	fmt.Println(req)
	url := "http://139.180.130.69"
	url = "http://45.32.112.148"
	resp, err := HttpRequest(url+"/getaddr", req, "6e347830f78d544jshsjld177f7278fb8")
	if err != nil {
		fmt.Println(err.Error())
		return
	}
	// CreateUsernameData 返回的数据 json
	type CreateUsernameData struct {
		Error   string `json:"error"`
		Message string `json:"message"`
		Result  string `json:"result"`
	}
	var dataAttr CreateUsernameData
	if err := json.Unmarshal([]byte(resp), &dataAttr); err != nil {
		return
	}
	fmt.Println(dataAttr)

}
