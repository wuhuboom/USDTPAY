package controller

import (
	"encoding/json"
	"example.com/m/dao/mysql"
	"example.com/m/dao/redis"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"io/ioutil"
	"math"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// ToAddress 地址管理
func ToAddress(c *gin.Context) {

	action := c.PostForm("action")
	//查看
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.ReceiveAddress, 0)
		Db := mysql.DB

		if add, isE := c.GetPostForm("address"); isE == true {
			Db = Db.Where("address=?", add)
		}

		if add, isE := c.GetPostForm("username"); isE == true {
			Db = Db.Where("username=?", add)
		}
		if add, isE := c.GetPostForm("kinds"); isE == true {
			Db = Db.Where("kinds=?", add)
		}

		if add, isE := c.GetPostForm("money"); isE == true {
			Db = Db.Where("money >=?", add)
		}
		//日期条件
		if start, isExist := c.GetPostForm("start_time"); isExist == true {
			if end, isExist := c.GetPostForm("end_time"); isExist == true {
				Db = Db.Where("updated >= ?", start).Where("updated<=?", end)
			}
		}
		var total int64
		Db.Table("receive_addresses").Count(&total)
		Db = Db.Model(&model.ReceiveAddress{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
	//查看地址账变
	if action == "getBalanceChange" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.AccountChange, 0)
		Db := mysql.DB
		var total int64
		//用户余额变动  receive_address_name
		if Rn, isE := c.GetQuery("receive_address_name"); isE == true {
			Db = Db.Where("receive_address_name=?", Rn)
		}
		Db.Table("account_changes").Count(&total)
		Db = Db.Model(&model.AccountChange{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}
		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
	//更新余额
	if action == "updateMoney" {
		re := make([]model.ReceiveAddress, 0)
		if Address, isE := c.GetPostForm("address"); isE == true {
			tools.ReturnError200(c, "执行成功")
			one := model.ReceiveAddress{}
			mysql.DB.Where("address=?", Address).First(&one)
			re = append(re, one)
		} else {
			result, _ := redis.Rdb.SetNX(c, "getAllAddressMoney", "getAllAddressMoney", time.Hour*5).Result()
			if result == false {
				tools.ReturnError200(c, "正在执行,不要重复执行")
				return
			}
			mysql.DB.Find(&re)
		}
		go func() {
			for _, v := range re {
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
					continue
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
			redis.Rdb.Del(c, "getAllAddressMoney")
			fmt.Println("检查成功!")
		}()
		tools.ReturnError200(c, "执行成功,等待结果")
		return
	}
	//资金归集
	if action == "getBalanceChange" {
		req := make(map[string]interface{})
		req["gas"] = c.PostForm("gas")
		req["min"] = c.PostForm("min")
		if req["gas"] == "" || req["min"] == "" {
			tools.ReturnError101(c, "非法参数")
			return
		}
		if addr, isExits := c.GetPostForm("addr"); isExits == true {
			if addr != "" {
				addArray := strings.Split(addr, "@")
				req["addrs"] = addArray
			}
		}
		req["ts"] = time.Now().UnixMilli()
		_, err := tools.HttpRequest(viper.GetString("project.ThreeUrl")+"/collect", req, viper.GetString("project.ApiKey"))
		if err != nil {
			tools.ReturnError101(c, "归集失败")
			log := model.Log{Kinds: 4, Ips: c.ClientIP(), Content: "资金归集失败,err:" + err.Error()}
			log.CreatedLogs(mysql.DB)
			return
		}
		log := model.Log{Kinds: 4, Ips: c.ClientIP(), Content: "资金归集成功"}
		log.CreatedLogs(mysql.DB)
		tools.ReturnError200(c, "归集成功")
		return
	}

}
