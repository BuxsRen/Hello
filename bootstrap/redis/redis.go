package redis

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"Hello/bootstrap/helper"
	"context"
	"fmt"
	"github.com/go-redis/redis/v8"
	"strconv"
)

var ctx = context.Background()
var Rdb *redis.Client

func init() {
	Rdb = connect()
}

// 连接Redis，返回实例
func connect() *redis.Client {
	var h = helper.Helper{}
	rdb := redis.NewClient(&redis.Options{
		Addr:     config.App.Redis.Host + ":" + config.App.Redis.Port,
		Password: config.App.Redis.PassWord,
		DB:       0,
	})
	e := rdb.Get(ctx, strconv.FormatInt(utils.GetTime(), 10)).Err()
	if e.Error() != "redis: nil" {
		fmt.Println("➦ " + e.Error())
		h.Exit("✘ Redis Connection Failed !", 3)
	}
	fmt.Printf("\U0001F9F1 Redis -> %v:%v\n", config.App.Redis.Host, config.App.Redis.Port)
	return rdb
}

