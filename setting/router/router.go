package router

import (
	"example.com/m/controller"
	"example.com/m/controller/h5"
	"example.com/m/controller/three"
	eeor "example.com/m/error"
	"fmt"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
)

func Setup() *gin.Engine {
	gin.SetMode(gin.ReleaseMode)
	r := gin.New()
	r.Use(Cors())
	r.Use(eeor.ErrHandler())
	r.NoMethod(eeor.HandleNotFound)
	r.NoRoute(eeor.HandleNotFound)
	//注册静态文件
	r.Static("/static", "./static")
	GroupV2 := r.Group("/v2").Use(PermissionToCheck())
	{
		{
			//管理员登录接口
			GroupV2.GET("/login", controller.Login)
			//管理员获取菜单  GetMenus
			GroupV2.GET("/getMenus", controller.GetMenus)
			GroupV2.GET("/getAllMoney", controller.GetAllMoney)

		}
		{
			//  拉起USDT订单
			GroupV2.POST("/createPrepaidPhoneOrders", h5.CreatePrepaidPhoneOrders)
			//大神接口回调
			GroupV2.POST("/getPayInformation", three.GetPayInformationBack)
		}
	}
	pageFist := r.Group("controller").Use(PermissionToCheck()).Use(PermissionToCheck())
	{
		pageFist.POST("fistPage", controller.ConsoleManagement)
	}
	//系统管理
	system := r.Group("system").Use(PermissionToCheck()).Use(PathUrlToCheck())
	{
		system.POST("parameterSetting", controller.Config)
		//GetRole  获取角色
		system.POST("roleManagement", controller.GetRole)
		//获取用户
		system.POST("userManagement", controller.GetAdmins)
	}

	//订单管理
	order := r.Group("order").Use(PermissionToCheck()).Use(PathUrlToCheck())
	{
		order.POST("topUpOrder", controller.TopUpOrder)
	}

	//地址管理
	address := r.Group("address").Use(PermissionToCheck()).Use(PathUrlToCheck())
	{
		address.POST("toAddress", controller.ToAddress)
	}

	//日志系统  LogBackManagement

	log := r.Group("log").Use(PermissionToCheck()).Use(PathUrlToCheck())
	{
		//backLog 回调日志
		log.POST("backLog", controller.LogBackManagement)
		//LogManagement  系统日志
		log.POST("systemLog", controller.LogManagement)

	}
	r.Run(fmt.Sprintf(":%d", viper.GetInt("project.port")))
	return r
}
