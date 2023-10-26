package controller

import (
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
	"strings"
	"time"
)

// GetRole 获取角色
func GetRole(c *gin.Context) {
	action := c.PostForm("action")
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.Role, 0)
		Db := mysql.DB
		//类型
		var total int64
		Db.Table("roles").Count(&total)
		Db = Db.Model(&model.Role{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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
	//新增角色
	if action == "add" {
		name := c.PostForm("name")
		if name == "" {
			tools.ReturnError101(c, "name 不能为空")
			return
		}
		err := mysql.DB.Save(&model.Role{Name: name, Created: time.Now().Unix()}).Error
		if err != nil {
			tools.ReturnError101(c, err.Error())
			return
		}
		tools.ReturnError200(c, "添加成功")
		return
	}
	//删除角色
	if action == "delete" {
		roleId := c.PostForm("role_id")
		if roleId == "1" {
			tools.ReturnError101(c, "超级管理员不可以被删除")
			return
		}
		mysql.DB.Model(&model.Role{}).Where("id=?", roleId).Delete(&model.Role{})
		mysql.DB.Model(&model.RoleMenu{}).Where("role_id=?", roleId).Delete(&model.RoleMenu{})
		tools.ReturnError200(c, "删除成功")
		return
	}
	//角色权限查看
	if action == "viewPermissions" {
		roleId := c.PostForm("role_id")
		//roleM := make([]model.RoleMenu, 0)
		//mysql.DB.Where("role_id=?", roleId).Find(&roleM)
		//tools.ReturnError200Data(c, roleM, "OK")
		mm := make([]model.Menu, 0)
		err := mysql.DB.Raw("SELECT menus.id,menus.name,menus.belong,menus.path,menus.menu_kind,menus.action,"+
			"menus.sort, menus.created FROM menus  LEFT JOIN  role_menus   "+
			"ON  role_menus.menu_id=menus.id WHERE  menus.menu_kind=1 AND "+
			" role_menus.role_id=?  AND  menus.belong=0 ORDER BY menus.sort  ASC", roleId).Scan(&mm).Error
		if err != nil {
			tools.ReturnError101(c, err.Error())
			return
		}
		RoleID, _ := strconv.Atoi(roleId)
		for i, menu := range mm {
			mm[i].Value = menu.ID

			m2 := make([]model.Menu, 0)
			mysql.DB.Where("belong=? and menu_kind=1 ", menu.ID).Find(&m2)
			for i2, m := range m2 {
				m2[i2].Value = m.ID
				per := make([]model.Menu, 0)
				mysql.DB.Where("belong=? and menu_kind=2", m.ID).Find(&per)
				m2[i2].Permissions = append(m2[i2].Permissions, per...)
				roleMenu := model.RoleMenu{MenuId: m.ID, RoleId: RoleID}
				if roleMenu.IfExist(mysql.DB) {
					m2[i2].IfChoose = true
					m2[i2].Check = true

				}
			}
			mm[i].SecondaryMenu = append(mm[i].SecondaryMenu, m2...)
			//添加permissions
			per := make([]model.Menu, 0)
			mysql.DB.Where("belong=? and menu_kind=2", menu.ID).Find(&per)

			for pi, p := range per {
				per[pi].Value = p.ID
				roleMenu := model.RoleMenu{MenuId: p.ID, RoleId: RoleID}
				if roleMenu.IfExist(mysql.DB) {
					per[pi].IfChoose = true
					per[pi].Check = true

				}
			}

			mm[i].Permissions = append(mm[i].Permissions, per...)

			//判断是否存在
			roleMenu := model.RoleMenu{MenuId: menu.ID, RoleId: RoleID}
			if roleMenu.IfExist(mysql.DB) {
				mm[i].IfChoose = true
				mm[i].Check = true
			}
		}
		tools.ReturnError200Data(c, mm, "OK")
		return
	}
	//更新角色权限
	if action == "updatePermissions" {
		roleId := c.PostForm("role_id")
		if roleId == "1" {
			tools.ReturnError101(c, "超级管理员权限不容修改")
			return
		}
		//角色不存在
		affected := mysql.DB.Where("id=?", roleId).Limit(1).Find(&model.Role{}).RowsAffected
		if affected == 0 {
			tools.ReturnError101(c, "修改的角色不存在")
			return
		}
		Men := c.PostForm("menu_id")
		MenArray := strings.Split(Men, "@")
		//先删除之前的所有权限
		err := mysql.DB.Model(&model.RoleMenu{}).Where("role_id=?", roleId).Delete(&model.RoleMenu{}).Error
		if err != nil {
			tools.ReturnError101(c, err.Error())
			return
		}
		for _, s := range MenArray {
			rowsAffected := mysql.DB.Where("id=?", s).Limit(1).Find(&model.Menu{}).RowsAffected
			if rowsAffected != 0 {
				roleIdInt, _ := strconv.Atoi(roleId)
				sInt, _ := strconv.Atoi(s)
				mysql.DB.Save(&model.RoleMenu{RoleId: roleIdInt, MenuId: sInt})
			}

		}
		tools.ReturnError200(c, "修改成功")
		return
	}

}

// GetAdmins 获取用户
func GetAdmins(c *gin.Context) {
	action := c.PostForm("action")
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.Admin, 0)
		Db := mysql.DB
		//类型
		var total int64
		Db.Table("admins").Count(&total)
		Db = Db.Model(&model.Admin{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
		err := Db.Find(&role).Error
		if err != nil {
			tools.ReturnError101(c, "ERR:"+err.Error())
			return
		}

		for i, admin := range role {
			m := model.Role{ID: admin.RoleId}
			role[i].RoleName = m.GetName(mysql.DB)
		}

		c.JSON(http.StatusOK, gin.H{
			"code":  0,
			"count": total,
			"data":  role,
		})
		return
	}
	//修改用户
	if action == "update" {

	}
	//添加用户
	if action == "add" {
		admin := model.Admin{}
		admin.Username = c.PostForm("username")
		admin.Password = tools.MD5(c.PostForm("password"))
		if c.PostForm("google_code") != "" {
			admin.GoogleCode = c.PostForm("google_code")
		}
		admin.RoleId, _ = strconv.Atoi(c.PostForm("role_id"))
		//判断角色是否存在
		affected := mysql.DB.Where("id=?", admin.RoleId).Limit(1).Find(&model.Role{}).RowsAffected
		if affected == 0 {
			tools.ReturnError101(c, "角色不存在")
			return
		}
		admin.Created = time.Now().Unix()
		admin.Token = string(tools.RandString(36))
		err := mysql.DB.Save(&admin).Error
		if err != nil {
			tools.ReturnError101(c, err.Error())
			return
		}
		tools.ReturnError200(c, "添加成功")
		return
	}

}
