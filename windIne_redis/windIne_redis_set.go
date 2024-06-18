package windIne_redis

import (
	"encoding/json"
	"fmt"
	"time"
)

// SAdd 集合--添加数据
func (gtr *WindIneRedis) SAdd(key1 string, values ...interface{}) error {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)

	err := gtr.redisClient.SAdd(ctx, aKey, values).Err()

	return err
}

// SMembers 集合--获取数据
func (gtr *WindIneRedis) SMembers(key1 string, values ...interface{}) []string {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)
	val, _ := gtr.redisClient.SMembers(ctx, aKey).Result()

	return val
}

// Scared 集合--获取数据数量
func (gtr *WindIneRedis) Scared(key1 string, key2 string) (Cnt int64) {
	aKey := fmt.Sprintf("%s:%s", prefix, key1)
	val, _ := gtr.redisClient.SCard(ctx, aKey).Result()

	return val
}

// SetMap 插入map
func (gtr *WindIneRedis) SetMap(key string, dictMap map[string]interface{}) error {
	jsonData, err := json.Marshal(dictMap)
	if err != nil {
		return err
	}
	return gtr.redisClient.Set(ctx, key, jsonData, 0).Err()
}

func (gtr *WindIneRedis) SetArrays(key string, dict []int64) error {
	jsonData, err := json.Marshal(dict)
	if err != nil {
		return err
	}
	return gtr.redisClient.Set(ctx, key, jsonData, 0).Err()
}

func (gtr *WindIneRedis) SetArraysExpire(key string, dict []int64, time time.Duration) error {
	jsonData, err := json.Marshal(dict)
	if err != nil {
		return err
	}
	return gtr.redisClient.Set(ctx, key, jsonData, time).Err()
}

func (gtr *WindIneRedis) SetMapExpire(key string, dictMap map[string]interface{}, time time.Duration) error {
	jsonData, err := json.Marshal(dictMap)
	if err != nil {
		return err
	}

	return gtr.redisClient.Set(ctx, key, jsonData, time).Err()
}

func (gtr *WindIneRedis) GetMap(key string) (map[string]interface{}, error) {
	jsonData, err := gtr.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var dictMap map[string]interface{}
	err = json.Unmarshal(jsonData, &dictMap)
	if err != nil {
		return nil, err
	}

	return dictMap, nil
}

func (gtr *WindIneRedis) GetArrays(key string) ([]int64, error) {
	jsonData, err := gtr.redisClient.Get(ctx, key).Bytes()
	if err != nil {
		return nil, err
	}

	var dict []int64
	err = json.Unmarshal(jsonData, &dict)
	if err != nil {
		return nil, err
	}
	return dict, nil
}
