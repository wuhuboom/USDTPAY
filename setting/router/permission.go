package router

import (
	"example.com/m/common"
	"example.com/m/dao/mysql"
	"example.com/m/dao/redis"
	"example.com/m/model"
	"example.com/m/tools"
	"fmt"
	"github.com/gin-gonic/gin"
)

// PermissionToCheck 权限校验  token 校验
func PermissionToCheck() gin.HandlerFunc {
	whiteUrl := []string{"/v2/login", "/v2/createPrepaidPhoneOrders", "/v2/getPayInformation", "/v2/getAllMoney"}
	return func(c *gin.Context) {
		if !tools.IsArray(whiteUrl, c.Request.URL.Path) {
			//token  校验
			//判断是用户还是管理员
			fmt.Println(c.Request.URL.Path)
			token := c.Request.Header.Get("token")
			//用户
			if len(token) == 36 {
				//管理员
				ad := model.Admin{}
				err := mysql.DB.Where("token=?", token).First(&ad).Error
				if err != nil {
					tools.JsonWrite(c, common.IllegalityCode, nil, "Sorry, your request is invalid")
					c.Abort()
					return
				}
				//判断token 是否过期?
				if redis.Rdb.Get(c, "AdminToken_"+token).Val() == "" {
					tools.JsonWrite(c, common.TokenExpire, nil, "Sorry, your login has expired")
					c.Abort()
					return
				}

				//设置who
				c.Set("who", ad)
				c.Next()
			} else {
				tools.JsonWrite(c, common.IllegalityCode, nil, "Sorry, your request is invalid")
				c.Abort()
				return
			}
		}
		c.Next()
	}

}

// PathUrlToCheck 路径权限控制

func PathUrlToCheck() gin.HandlerFunc {
	return func(c *gin.Context) {
		// 获取管理员
		who, _ := c.Get("who")
		Admin := who.(model.Admin)
		menu := model.Menu{}
		action := c.PostForm("action")
		if action == "" {
			tools.JsonWrite(c, common.IllegalityCode, nil, "Illegal request->>")
			c.Abort()
			return
		}
		affected := mysql.DB.Where("path=? and action=?", c.Request.URL.Path, action).Limit(1).Find(&menu).RowsAffected
		if affected == 0 {
			tools.JsonWrite(c, common.IllegalityCode, nil, "Illegal request-")
			c.Abort()
			return
		}
		rowsAffected := mysql.DB.Where("role_id=? and menu_id=?", Admin.ID, menu.ID).Limit(1).Find(&model.RoleMenu{}).RowsAffected
		if rowsAffected == 0 {
			tools.JsonWrite(c, common.IllegalityCode, nil, "Illegal request--")
			c.Abort()
			return
		}
		//fmt.Println(c.Request.URL.Path)
		c.Next()
	}

}
