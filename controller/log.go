package controller

import (
	"example.com/m/dao/mysql"
	"example.com/m/model"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
	"net/http"
	"strconv"
)

// LogManagement 系统日志管理
func LogManagement(c *gin.Context) {
	action := c.PostForm("action")
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.Log, 0)
		Db := mysql.DB
		//类型
		if content, isExist := c.GetPostForm("kinds"); isExist == true {
			Db = Db.Where("kinds=?", content)
		}
		//日期条件
		if start, isExist := c.GetPostForm("start_time"); isExist == true {
			if end, isExist := c.GetPostForm("end_time"); isExist == true {
				Db = Db.Where("created >= ?", start).Where("created<=?", end)
			}
		}
		var total int64
		Db.Table("logs").Count(&total)
		Db = Db.Model(&model.Log{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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

}

// LogBackManagement 回调日志管理
func LogBackManagement(c *gin.Context) {

	action := c.PostForm("action")
	if action == "check" {
		page, _ := strconv.Atoi(c.PostForm("page"))
		limit, _ := strconv.Atoi(c.PostForm("limit"))
		role := make([]model.BackLog, 0)

		Db := mysql.DB
		//类型
		if content, isExist := c.GetPostForm("kinds"); isExist == true {
			Db = Db.Where("kinds=?", content)
		}

		//TxHash
		if content, isExist := c.GetPostForm("tx_hash"); isExist == true {
			Db = Db.Where("tx_hash=?", content)
		}

		//json_content
		if content, isExist := c.GetPostForm("json_content"); isExist == true {
			Db = Db.Where("json_content like  ?", "%"+content+"%")
		}

		//日期条件
		if start, isExist := c.GetPostForm("start_time"); isExist == true {
			if end, isExist := c.GetPostForm("end_time"); isExist == true {
				Db = Db.Where("created >= ?", start).Where("created<=?", end)
			}
		}

		var total int64
		Db.Table("back_logs").Count(&total)
		Db = Db.Model(&model.BackLog{}).Offset((page - 1) * limit).Limit(limit).Order("created desc")
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

}
