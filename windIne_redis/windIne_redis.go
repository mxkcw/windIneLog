package windIne_redis

import (
	"context"
	"fmt"
	"github.com/mxkcw/windIneLog/windIne_log"
	"github.com/redis/go-redis/v9"
	"time"
)

var (
	isSetup    bool
	OwnerRedis *WindIneRedis
	ctx        = context.Background()
	prefix     string
)

type RedisConfig struct {
	Addr       string `yaml:"addr" json:"addr"` // address
	Port       string `yaml:"port" json:"port"`
	Pwd        string `yaml:"pwd" json:"pwd"`                // pwd
	SocketBuck int    `yaml:"socketBuck" json:"socket_buck"` // 插槽
}

type WindIneRedis struct {
	cfg         *RedisConfig
	redisClient *redis.Client
}

// SetupRedisConnection 初始化Redis连接
func SetupRedisConnection(redisCfg RedisConfig, prefixStr string) (success bool) {
	if isSetup == false {
		OwnerRedis = &WindIneRedis{
			cfg: &redisCfg,
			redisClient: redis.NewClient(&redis.Options{
				Addr:     fmt.Sprintf("%s:%s", redisCfg.Addr, redisCfg.Port),
				Password: redisCfg.Pwd,        // no password set
				DB:       redisCfg.SocketBuck, // use default DB
			}),
		}

		prefix = prefixStr
		err := OwnerRedis.redisClient.Set(ctx, "hello", "helloValue", 0).Err()
		if err != nil {
			windIne_log.LogErrorf("[redis setup] error [%s]", err)
			isSetup = false
		} else {
			windIne_log.LogInfof("[redis setup] [%s]", "success")
			isSetup = true
		}
	}
	return isSetup
}

// Set 插入单条数据
func (gtr *WindIneRedis) Set(key string, value string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	err := gtr.redisClient.Set(ctx, aKey, value, 0).Err()

	return err
}

// SetStrByExpire 插入过期时间
func (gtr *WindIneRedis) SetStrByExpire(key, value string, time time.Duration) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	err := gtr.redisClient.Set(ctx, aKey, value, time).Err()
	return err
}

// Get 获取单条数据
func (gtr *WindIneRedis) Get(key string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := gtr.redisClient.Get(ctx, aKey).Result()

	return val, err
}

// Del 删除单条数据
func (gtr *WindIneRedis) Del(key string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := gtr.redisClient.Del(ctx, aKey).Err()

	return err
}

// Keys 删除单条数据
func (gtr *WindIneRedis) Keys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := gtr.redisClient.Keys(ctx, aKey).Result()

	return val, err
}

// Exists 检测是否过期
func (gtr *WindIneRedis) Exists(key string) bool {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, _ := gtr.redisClient.Exists(ctx, aKey).Result()
	if val > 0 {
		return true
	} else {
		return false
	}
}
