package model

import (
	"fmt"
	"gorm.io/gorm"
	"time"
)

type Menu struct {
	ID            int    `json:"id" gorm:"primaryKey;comment:'主键'"`
	Name          string `json:"name" gorm:"uniqueIndex"`
	Belong        int    `json:"belong" `                    //0顶级菜单
	Path          string `json:"path"`                       //路劲配置
	MenuKind      int    `json:"menu_kind" gorm:"default:1"` //菜单类型 1 菜单 2功能
	Action        string `json:"action"`                     //增删改查
	Sort          int    `json:"sort"`                       //排序 小的在上面
	Created       int64  `json:"created"`
	SecondaryMenu []Menu `json:"secondary_menu"  gorm:"-"`
	Permissions   []Menu `json:"permissions"  gorm:"-"`
}

// CheckIsExistModelMenu CheckIsExistModelAdmin 创建
func CheckIsExistModelMenu(db *gorm.DB) {
	if db.Migrator().HasTable(&Menu{}) {
		fmt.Println("数据库已经存在了!")
		db.AutoMigrate(&Menu{})
	} else {
		fmt.Println("数据不存在,所以我要先创建数据库")
		db.Migrator().CreateTable(&Menu{})
		me := make([]Menu, 0)
		me = append(me, Menu{ID: 1, Name: "控制台", Belong: 0, Created: time.Now().Unix(), Sort: 1, Path: "controller"})
		//me = append(me, Menu{ID: 3, Name: "日志管理", Belong: 0, Created: time.Now().Unix(), Sort: 2})

		{
			me = append(me, Menu{ID: 1111, Name: "查看", Belong: 1, Created: time.Now().Unix(), Path: "/controller/fistPage", MenuKind: 2, Action: "check"})
		}
		me = append(me, Menu{ID: 5, Name: "订单管理", Belong: 0, Created: time.Now().Unix(), Sort: 3, Path: "order"})
		{
			me = append(me, Menu{ID: 510, Name: "充值订单", Belong: 5, Created: time.Now().Unix(), Path: "/order/topUpOrder"})
			{
				me = append(me, Menu{ID: 5101, Name: "查看", Belong: 510, Created: time.Now().Unix(), Path: "/order/topUpOrder", MenuKind: 2, Action: "check"})
				me = append(me, Menu{ID: 5102, Name: "修改", Belong: 510, Created: time.Now().Unix(), Path: "/order/topUpOrder", MenuKind: 2, Action: "update"})
				me = append(me, Menu{ID: 5103, Name: "回调订单", Belong: 510, Created: time.Now().Unix(), Path: "/order/orderBack", MenuKind: 2, Action: "orderBack"})

			}
		}

		me = append(me, Menu{ID: 6, Name: "地址管理", Belong: 0, Created: time.Now().Unix(), Sort: 4, Path: "address"})

		{
			me = append(me, Menu{ID: 610, Name: "收账地址", Belong: 6, Created: time.Now().Unix(), Path: "/address/toAddress"})
			{
				me = append(me, Menu{ID: 6101, Name: "查看", Belong: 610, Created: time.Now().Unix(), Path: "/address/toAddress", MenuKind: 2, Action: "check"})
				me = append(me, Menu{ID: 6102, Name: "查看账变", Belong: 610, Created: time.Now().Unix(), Path: "/address/toAddress", MenuKind: 2, Action: "getBalanceChange"})
				me = append(me, Menu{ID: 6103, Name: "更新余额", Belong: 610, Created: time.Now().Unix(), Path: "/address/toAddress", MenuKind: 2, Action: "updateMoney"})
				me = append(me, Menu{ID: 6104, Name: "资金归集", Belong: 610, Created: time.Now().Unix(), Path: "/address/toAddress", MenuKind: 2, Action: "collectByYourself"})

			}
		}

		me = append(me, Menu{ID: 7, Name: "日志管理", Belong: 0, Created: time.Now().Unix(), Sort: 4, Path: "log"})
		{
			me = append(me, Menu{ID: 710, Name: "回调日志", Belong: 7, Created: time.Now().Unix(), Path: "/log/backLog"})
			{
				me = append(me, Menu{ID: 7101, Name: "查看", Belong: 710, Created: time.Now().Unix(), Path: "/log/backLog", MenuKind: 2, Action: "check"})

			}
			me = append(me, Menu{ID: 711, Name: "系统日志", Belong: 7, Created: time.Now().Unix(), Path: "/log/systemLog"})
			{
				me = append(me, Menu{ID: 7102, Name: "查看", Belong: 711, Created: time.Now().Unix(), Path: "/log/systemLog", MenuKind: 2, Action: "check"})

			}
		}

		me = append(me, Menu{ID: 4, Name: "系统管理", Belong: 0, Created: time.Now().Unix(), Path: "/system", Sort: 8})
		{
			me = append(me, Menu{ID: 410, Name: "系统参数", Belong: 4, Created: time.Now().Unix(), Path: "/system/parameterSetting"})
			{
				me = append(me, Menu{ID: 4101, Name: "查看", Belong: 410, Created: time.Now().Unix(), Path: "/system/parameterSetting", MenuKind: 2, Action: "check"})
				me = append(me, Menu{ID: 4102, Name: "修改", Belong: 410, Created: time.Now().Unix(), Path: "/system/parameterSetting", MenuKind: 2, Action: "update"})
			}
			me = append(me, Menu{ID: 411, Name: "用户管理", Belong: 4, Created: time.Now().Unix(), Path: "/system/userManagement"})

			{
				me = append(me, Menu{ID: 4111, Name: "查看", Belong: 411, Created: time.Now().Unix(), Path: "/system/userManagement", MenuKind: 2, Action: "check"})
				me = append(me, Menu{ID: 4112, Name: "修改", Belong: 411, Created: time.Now().Unix(), Path: "/system/userManagement", MenuKind: 2, Action: "update"})
				me = append(me, Menu{ID: 4113, Name: "删除", Belong: 411, Created: time.Now().Unix(), Path: "/system/userManagement", MenuKind: 2, Action: "delete"})
				me = append(me, Menu{ID: 4114, Name: "添加", Belong: 411, Created: time.Now().Unix(), Path: "/system/userManagement", MenuKind: 2, Action: "add"})
			}

			me = append(me, Menu{ID: 412, Name: "角色管理", Belong: 4, Created: time.Now().Unix(), Path: "/system/roleManagement"})
			{
				me = append(me, Menu{ID: 4121, Name: "查看", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "check"})
				me = append(me, Menu{ID: 4122, Name: "修改", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "update"})
				me = append(me, Menu{ID: 4123, Name: "删除", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "delete"})
				me = append(me, Menu{ID: 4124, Name: "添加", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "add"})
				me = append(me, Menu{ID: 4125, Name: "权限查看", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "viewPermissions"})
				me = append(me, Menu{ID: 4126, Name: "修改权限", Belong: 412, Created: time.Now().Unix(), Path: "/system/roleManagement", MenuKind: 2, Action: "updatePermissions"})

			}
		}

		db.Create(me)
	}
}
