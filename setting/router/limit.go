package router

import (
	"example.com/m/common"
	"example.com/m/dao/redis"
	"example.com/m/tools"
	"github.com/gin-gonic/gin"
	"github.com/spf13/viper"
	"time"
)

func LimitIpRequestSameUrlForUser() gin.HandlerFunc {
	return func(context *gin.Context) {
		//获取请求的地址
		urlPath := context.Request.URL.Path
		ip := context.ClientIP()
		//回去系统设置的ip限制次数
		var LimitTimes int
		LimitTimes = viper.GetInt("project.limitIp")
		key := ip + urlPath
		curr := redis.Rdb.LLen(context, key).Val()
		if int(curr) >= LimitTimes {
			//超出了限制
			tools.JsonWrite(context, common.IpLimitWaring, nil, "Don't ask too fast")
			context.Abort()
		}
		if v := redis.Rdb.Exists(context, key).Val(); v == 0 {
			pipe := redis.Rdb.TxPipeline()
			pipe.RPush(context, key, key)
			//设置过期时间
			pipe.Expire(context, key, 1*time.Second)
			_, _ = pipe.Exec(context)
		} else {
			redis.Rdb.RPushX(context, key, key)
		}
		context.Next()
	}
}
