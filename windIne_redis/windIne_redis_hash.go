package windIne_redis

import (
	"fmt"
)

func (gtr *WindIneRedis) HGetAll(key string) (map[string]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	val, err := gtr.redisClient.HGetAll(ctx, aKey).Result()
	return val, err
}

// HSet Hash类型-插入单条数据
func (gtr *WindIneRedis) HSet(key string, subKey string, jsonByte []byte) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := gtr.redisClient.HSet(ctx, aKey, subKey, jsonByte).Err()
	return err
}

// HGet Hash类型-获取单条数据
func (gtr *WindIneRedis) HGet(key string, subKey string) (string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, aErr := gtr.redisClient.HGet(ctx, aKey, subKey).Result()
	return val, aErr
}

// HDel Hash类型-删除单条数据
func (gtr *WindIneRedis) HDel(key string, subKey string) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key)

	err := gtr.redisClient.HDel(ctx, aKey, subKey).Err()
	return err
}

// HExists Hash类型-判断是否存在
func (gtr *WindIneRedis) HExists(key string, subKey string) bool {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, _ := gtr.redisClient.HExists(ctx, aKey, subKey).Result()

	return val
}

func (gtr *WindIneRedis) HKeys(key string) ([]string, error) {
	aKey := fmt.Sprintf("%s:%s", prefix, key)
	val, err := gtr.redisClient.HKeys(ctx, aKey).Result()

	return val, err
}
