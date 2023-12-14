package process

import (
	"fmt"
	"math/big"
	"testing"
)

func TestGetNewBlock(t *testing.T) {
	//TCdxXudGQNeDTVj8qp1X1jNkJeepVdgkey  我原地址是这个
	//0000000000000000000000411d44658e2acbaf0b1f0c2ccb69904f470b6b27a1转成了这样子
	//base58Address := "TCdxXudGQNeDTVj8qp1X1jNkJeepVdgkey"
	//
	//hexAddress, err := base58ToHex(base58Address)
	//if err != nil {
	//	fmt.Println("Error:", err)
	//	return
	//}
	//
	//fmt.Println("Hex Address:", hexAddress)

	hexString := "000000000000000000000000000000000000000000000000000000002f6355de"

	decimal, err := toDecimal(hexString)
	if err != nil {
		fmt.Println("Failed to convert hex to decimal:", err)
		return
	}

	fmt.Println(decimal.String())
	divisor := big.NewInt(1000000)
	pp := big.NewInt(795039198)
	fmt.Println(new(big.Int).Div(pp, divisor))
}

//func toDecimal(hexString string) (*big.Int, error) {
//	hexString = strings.TrimPrefix(hexString, "0x")
//	decimal := new(big.Int)
//	_, success := decimal.SetString(hexString, 16)
//	if !success {
//		return nil, fmt.Errorf("Failed to convert hex string to decimal")
//	}
//	return decimal, nil
//}
