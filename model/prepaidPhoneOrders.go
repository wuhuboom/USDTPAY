package model

import (
	"encoding/base64"
	"encoding/json"
	"example.com/m/util"
	"fmt"
	"gorm.io/gorm"
	"math/rand"
	"strconv"
	"time"
)

// PrepaidPhoneOrders 充值订单
type PrepaidPhoneOrders struct {
	ID            int    `gorm:"primaryKey;comment:'主键'" json:"id"`
	PlatformOrder string `json:"platform_order"` //平台订单  (前台订单号)
	ThreeOrder    string `json:"three_order"`    //三方订单  (自己 随机成功)
	//ToAddress        string  `json:"to_address"`                                  //收账地址(庄家地址)
	//FormAddress      string  `json:"form_address"`                                //转账地址(玩家地址)
	RechargeAddress   string  `json:"recharge_address"`                            //平台地址
	CollectionAddress string  `json:"collection_address"`                          //玩家地址
	Username          string  `json:"username"`                                    //充值用户名
	AccountOrders     float64 `gorm:"type:decimal(10,2)" json:"account_orders"`    //充值金额 (订单金额)
	AccountPractical  float64 `gorm:"type:decimal(10,2)" json:"account_practical"` //充值金额(实际返回金额)
	Status            int     `json:"status"`                                      //订单状态  1 未支付  2已经支付了  3已经失效
	ThreeBack         int     `json:"three_back" gorm:"default:1"`                 //三方回调 1未回调  2已结回调
	Created           int64   `json:"created"`                                     //订单创建时间
	Updated           int64   `json:"updated"`                                     //更新时间(回调时间)
	Successfully      int64   `json:"successfully"`                                //交易成功 时间(区块时间戳)
	Date              string  `json:"date"`                                        //日期
	BackUrl           string  `json:"back_url"`                                    //回调的地址
	Remark            string  `json:"remark"`                                      //备注
	RechargeType      string  `json:"recharge_type"`                               //充值类型
	TxHash            string  `json:"tx_hash"`                                     //回调的hash值
	BackData          string  `gorm:"type:text" json:"back_data"`                  //回调数据
	FootballBackData  string  `json:"football_back_data" gorm:"type:text"`         //返回数据
	ErrString         string  `json:"err_string"`                                  //错误信息
}

func CheckIsExistModelPrepaidPhoneOrders(db *gorm.DB) {
	if db.Migrator().HasTable(&PrepaidPhoneOrders{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&PrepaidPhoneOrders{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&PrepaidPhoneOrders{})
	}
}

// IfUseThisTxHash 是否存在这个Hash 值
func (p *PrepaidPhoneOrders) IfUseThisTxHash(db *gorm.DB) bool {
	affected := db.Where("tx_hash=?", p.TxHash).Limit(1).Find(&PrepaidPhoneOrders{}).RowsAffected
	if affected == 0 {
		return false
	}
	return true
}

// UpdateMaxCreatedOfStatusToTwo 专属地址回调
func (p *PrepaidPhoneOrders) UpdateMaxCreatedOfStatusToTwo(db *gorm.DB, OrderEffectivityTime int64) bool {
	type Create struct {
		PlatformOrder    string
		RechargeAddress  string
		Username         string
		AccountOrders    float64 //订单充值金额
		AccountPractical float64 //  实际充值的金额
		RechargeType     string
		BackUrl          string
	}
	//找到这条数据
	pp := PrepaidPhoneOrders{}
	err := db.Where("username=?", p.Username).Where("status= ? and recharge_type= ?", 1, p.RechargeType).Last(&pp).Error
	if err == nil {
		if time.Now().Unix()-pp.Created <= OrderEffectivityTime {
			//找到最新的数据(并且在有效时间累)
			//这里 要回调给前台
			update := PrepaidPhoneOrders{
				Updated:           time.Now().Unix(),
				Successfully:      p.Successfully,
				ThreeBack:         2,
				Status:            2,
				AccountPractical:  p.AccountPractical,
				CollectionAddress: p.CollectionAddress, TxHash: p.TxHash}
			if pp.BackUrl != "" {
				var tt Create
				tt.PlatformOrder = pp.PlatformOrder
				tt.RechargeAddress = p.RechargeAddress
				tt.Username = p.Username
				tt.AccountOrders = pp.AccountOrders
				tt.AccountPractical = p.AccountPractical
				tt.RechargeType = p.RechargeType
				data, err := json.Marshal(tt)
				if err != nil {
					return false
				}
				//更新
				update.BackData = string(data)
				data, err = util.RsaEncryptForEveryOne(data)
				util.BackUrlToPay(pp.BackUrl, base64.StdEncoding.EncodeToString(data))
			}
			db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Updates(&update)
			return true
		} else {
			db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Updates(&PrepaidPhoneOrders{Status: 3, Updated: time.Now().Unix()})
			return false
		}
	}

	//没有这条数据  默认是线下支付   自行创建这笔订单 Username: p.UserID, Successfully: p.Timestamp, AccountPractical: p.Amount}
	pt := PrepaidPhoneOrders{}
	pt.Created = time.Now().Unix()
	pt.Updated = 0
	pt.ThreeOrder = time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000))
	pt.PlatformOrder = ""
	pt.Status = 2
	pt.ThreeBack = 2
	pt.Username = p.Username
	pt.Successfully = p.Successfully
	pt.AccountPractical = p.AccountPractical
	pt.RechargeType = p.RechargeType
	pt.RechargeAddress = p.RechargeAddress
	pt.CollectionAddress = p.CollectionAddress
	pt.Date = time.Now().Format("2006-01-02")
	pt.TxHash = p.TxHash
	db.Save(&pt)
	return false
}

// UpdatePondOrderCratedAndUpdated 池地址回调
func (p *PrepaidPhoneOrders) UpdatePondOrderCratedAndUpdated(db *gorm.DB) bool {
	type Create struct {
		PlatformOrder    string
		RechargeAddress  string
		Username         string
		AccountOrders    float64 //订单充值金额
		AccountPractical float64 //  实际充值的金额
		RechargeType     string
		BackUrl          string
	}
	//找到这条数据
	pp := PrepaidPhoneOrders{}
	config := Config{}
	config.Expiration = 30
	db.Where("id=?", 1).First(&config)
	//	创建+过期  > 现在时间
	err := db.Where("recharge_address=?", p.RechargeAddress).
		Where("status= ? and recharge_type= ? and created > ?", 1, p.RechargeType, time.Now().Unix()-config.Expiration*60).
		First(&pp).Error
	if err == nil {
		//找到了这笔订单
		updateData := PrepaidPhoneOrders{
			Updated:           time.Now().Unix(),
			Successfully:      p.Successfully,
			ThreeBack:         2,
			Status:            2,
			AccountPractical:  p.AccountPractical,
			CollectionAddress: p.CollectionAddress,
			Date:              time.Now().Format("2006-01-02"),
			TxHash:            p.TxHash,
		}

		//这里 要回调给前台
		if pp.BackUrl != "" {
			var tt Create
			tt.PlatformOrder = pp.PlatformOrder
			tt.RechargeAddress = p.RechargeAddress
			tt.Username = p.Username
			tt.AccountOrders = pp.AccountOrders
			tt.AccountPractical = p.AccountPractical
			tt.RechargeType = p.RechargeType
			data, err := json.Marshal(tt)
			if err != nil {
				return false
			}
			updateData.BackData = string(data)
			data, err = util.RsaEncryptForEveryOne(data)
			util.BackUrlToPay(pp.BackUrl, base64.StdEncoding.EncodeToString(data))
		}
		db.Model(&PrepaidPhoneOrders{}).Where("id=?", pp.ID).Updates(&updateData)
		return true
	}

	//没找到这个订单 补这个订单
	//没有这条数据  默认是线下支付   自行创建这笔订单 Username: p.UserID, Successfully: p.Timestamp, AccountPractical: p.Amount}
	pt := PrepaidPhoneOrders{}
	pt.Created = time.Now().Unix()
	pt.Updated = 0
	pt.ThreeOrder = time.Now().Format("20060102150405") + strconv.Itoa(rand.Intn(100000))
	pt.PlatformOrder = ""
	pt.Status = 2
	pt.ThreeBack = 2
	pt.Username = p.Username
	pt.Successfully = p.Successfully
	pt.AccountPractical = p.AccountPractical
	pt.RechargeType = p.RechargeType
	pt.RechargeAddress = p.RechargeAddress
	pt.CollectionAddress = p.CollectionAddress
	pt.Date = time.Now().Format("2006-01-02")
	pt.TxHash = p.TxHash
	db.Save(&pt)
	return true
}

func (p *PrepaidPhoneOrders) HashExist(db *gorm.DB) bool {
	affected := db.Where("tx_hash=?", p.TxHash).Limit(1).Find(p).RowsAffected
	if affected == 0 {
		return false
	}
	return true
}
