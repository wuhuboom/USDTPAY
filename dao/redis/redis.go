package redis

import (
	"context"
	"fmt"
	"github.com/redis/go-redis/v9"
	"github.com/spf13/viper"
)

var Rdb *redis.Client

//var LuaRdb *redis.Script
//var RdbReader *redis.Client

func Init() (err error) {
	Rdb = redis.NewClient(&redis.Options{
		Addr: fmt.Sprintf("%s:%d",
			viper.GetString("redis.host"),
			viper.GetInt("redis.port"),
		),
		Password: viper.GetString("redis.password"),
		DB:       viper.GetInt("redis.db"),
		PoolSize: viper.GetInt("redis.pool_size"),
	})

	ctx := context.Background()
	result, err := Rdb.Ping(ctx).Result()

	fmt.Println(result)
	return err
}
