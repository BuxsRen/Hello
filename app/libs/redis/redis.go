package redis

import (
	"Hello/app/libs/utils"
	"Hello/bootstrap/config"
	"Hello/bootstrap/redis"
	"context"
	"log"
	"time"
)

var rdb = *redis.Rdb
var ctx = context.Background()
var prefix = config.App.Redis.Prefix

// Redis 封装库
/**
 *@Example:
	rdb := redis.Redis{}
	fmt.Println(rdb.Get("test"))
*/
type Redis struct{}

type Api interface {
	Get(key string) string
	Set(key string, val interface{})
	Setex(key string, val interface{}, t int)
	Del(key string)
	RincrBy(key string, mode string)
}

// 取值。键
func (this *Redis) Get(key string) string {
	val, e := rdb.Get(ctx, prefix+key).Result()
	if e != nil {
		return ""
	}
	return val
}

// 存值。键，值
func (this *Redis) Set(key string, val interface{}) {
	e := rdb.Set(ctx, prefix+key, val, 0).Err()
	if e != nil {
		log.Println(e)
	}
}

// 存值，带过期时间，单位秒。键，值，秒
func (this *Redis) Setex(key string, val interface{}, t int) {
	s := time.Duration(t) * time.Second
	e := rdb.Set(ctx, prefix+key, val, s).Err()
	if e != nil {
		log.Println(e)
	}
}

// 删除值。键
func (this *Redis) Del(key string) {
	e := rdb.Del(ctx, prefix+key).Err()
	if e != nil {
		log.Println(e)
	}
}

// 递增或递减。键，模式(+,-)
func (this *Redis) RincrBy(key string, mode string) {
	if mode == "+" {
		rdb.IncrBy(ctx, prefix+key, 1)
	} else {
		rdb.IncrBy(ctx, prefix+key, -1)
	}
}

// 输出缓存信息
func (this *Redis) CheckOut(key string) {
	str := this.Get(key)
	if str != "" {
		utils.ExitJson(str)
	}
}