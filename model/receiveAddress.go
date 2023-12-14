package model

import (
	"encoding/json"
	"example.com/m/tools"
	"fmt"
	"github.com/spf13/viper"
	"go.uber.org/zap"
	"gorm.io/gorm"
	"math"
	"time"
)

type ReceiveAddress struct {
	ID                  int     `gorm:"primaryKey;comment:'主键'" json:"id"`
	Username            string  `gorm:"uniqueIndex" json:"username"`
	ReceiveNums         int     `json:"receive_nums"`                                //收款笔数
	LastGetAccount      float64 `gorm:"type:decimal(10,2)"  json:"last_get_account"` //最后一次的入账金额
	Address             string  `gorm:"uniqueIndex" json:"address"`                  //收账地址
	Money               float64 `gorm:"type:decimal(10,2)" json:"money"`             //账户余额
	TheLastGetMoneyTime int64   `gorm:"default:0" json:"the_last_get_money_time"`    //最后一次获取余额的时间
	Kinds               int     `gorm:"default:1" json:"kinds"`                      //地址类型  1普通玩家地址  2 池地址
	LastUseTime         int64   `gorm:"default:0" json:"lastUseTime"`                // 最后一次使用时间
	Status              int     `gorm:"default:1" json:"status"`                     // 1状态正常 2   关闭
	Created             int64   `json:"created"`
	Updated             int64   `json:"updated"`
}

func (r *ReceiveAddress) IsExist(db *gorm.DB) (bool, *ReceiveAddress) {
	affected := db.Where("address=?", r.Address).Limit(1).Find(r).RowsAffected
	if affected > 0 {
		return true, r
	}
	return false, r
}

func CheckIsExistModelReceiveAddress(db *gorm.DB) {
	if db.Migrator().HasTable(&ReceiveAddress{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&ReceiveAddress{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&ReceiveAddress{})
	}
}

// CreateUsername 创建这个用户  获取用户收款地址
func (r *ReceiveAddress) CreateUsername(db *gorm.DB, url string) ReceiveAddress {
	r.Created = time.Now().Unix()
	r.ReceiveNums = 0
	r.LastGetAccount = 0
	//获取收账地址  url 请求  {"error":"0","message":"","result":"4564554545454545"}   //返回数据
	req := make(map[string]interface{})
	req["user"] = r.Username
	req["ts"] = time.Now().UnixMilli()
	resp, err := tools.HttpRequest(url+"/getaddr", req, viper.GetString("project.ApiKey"))
	if err != nil {
		fmt.Println(err.Error())
		return ReceiveAddress{}
	}
	// CreateUsernameData 返回的数据 json
	type CreateUsernameData struct {
		Error   int    `json:"error"`
		Message string `json:"message"`
		Result  string `json:"result"`
		User    string `json:"user"`
	}

	var dataAttr CreateUsernameData
	if err := json.Unmarshal([]byte(resp), &dataAttr); err != nil {
		fmt.Println(err)
		return ReceiveAddress{}
	}

	fmt.Println(dataAttr)

	if dataAttr.Result != "" {
		r.Address = dataAttr.Result
		err := db.Save(&r).Error
		if err != nil {
			return ReceiveAddress{}
		}
	}
	fmt.Println(r)
	return *r
}

// UpdateReceiveAddressLastInformationTo0 地址余额清零
func (r *ReceiveAddress) UpdateReceiveAddressLastInformationTo0(db *gorm.DB) bool {
	re := ReceiveAddress{}
	err := db.Where("username=?", r.Username).First(&re).Error
	if err == nil {
		zap.L().Debug("余额清0,用户:" + r.Username)
		updated := make(map[string]interface{})
		updated["Updated"] = r.Updated
		updated["Money"] = 0
		err := db.Model(&ReceiveAddress{}).Where("id=?", re.ID).Updates(updated).Error
		if err == nil {
			return true
		}
	}
	return false
}

// UpdateReceiveAddressLastInformation 更新钱包地址
func (r *ReceiveAddress) UpdateReceiveAddressLastInformation(db *gorm.DB) bool {
	re := ReceiveAddress{}
	err := db.Where("username=?", r.Username).First(&re).Error
	if err == nil {
		nums := re.ReceiveNums + 1
		err := db.Model(&ReceiveAddress{}).Where("id=?", re.ID).Updates(
			&ReceiveAddress{ReceiveNums: nums, LastGetAccount: r.LastGetAccount, Updated: r.Updated, Money: r.Money}).Error
		if err == nil {
			//更新账变
			change := AccountChange{
				ChangeAmount: math.Abs(re.Money - r.Money),
				Kinds:        2, OriginalAmount: re.Money,
				NowAmount:          r.Money,
				ReceiveAddressName: r.Username}
			change.Add(db)
			return true
		}
	}
	return false
}
