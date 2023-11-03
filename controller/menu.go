package controller

import (
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
)

// GetMenus 获取菜单
func GetMenus(c *gin.Context) {
	admin := GetWho(c)
	mm := make([]model.Menu, 0)
	err := mysql.DB.Raw("SELECT menus.id,menus.name,menus.belong,menus.path,menus.menu_kind,menus.action,"+
		"menus.sort, menus.created FROM menus  LEFT JOIN  role_menus   "+
		"ON  role_menus.menu_id=menus.id WHERE  menus.menu_kind=1 AND "+
		" role_menus.role_id=?  AND  menus.belong=0 ORDER BY menus.sort  ASC", admin.RoleId).Scan(&mm).Error
	if err != nil {
		tools.ReturnError101(c, err.Error())
		return
	}
	for i, menu := range mm {
		m2 := make([]model.Menu, 0)
		mysql.DB.Where("belong=? and menu_kind=1 ", menu.ID).Find(&m2)
		//mysql.DB.Raw("SELECT menus.id,menus.name,menus.belong,menus.path,menus.menu_kind,menus.action,menus.sort, menus.created FROM menus  LEFT JOIN  role_menus   "+
		//	"ON  role_menus.menu_id=menus.id WHERE  "+
		//	" role_menus.role_id=?  AND  menus.belong=?  ORDER BY menus.sort  ASC", admin.RoleId, menu.ID).Scan(&m2)

		for i2, m := range m2 {
			per := make([]model.Menu, 0)
			mysql.DB.Where("belong=? and menu_kind=2", m.ID).Find(&per)
			m2[i2].Permissions = append(m2[i2].Permissions, per...)
		}

		mm[i].SecondaryMenu = append(mm[i].SecondaryMenu, m2...)
		//添加permissions
		per := make([]model.Menu, 0)
		mysql.DB.Where("belong=? and menu_kind=2", menu.ID).Find(&per)
		mm[i].Permissions = append(mm[i].Permissions, per...)

	}
	tools.ReturnError200Data(c, mm, "OK")
	return
}

// GetWho 获取自己
func GetWho(c *gin.Context) model.Admin {
	who, _ := c.Get("who")
	return who.(model.Admin)
}
