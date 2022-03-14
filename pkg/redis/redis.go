package redis

import (
	"context"

	"github.com/go-redis/redis/v8"
)

type IRedisClient interface {
	ZRangeAdd(listKey, key string)
	ZRange(listKey string) []redis.Z
	ZRangeScore(listKey, key string) int64
	ZIncrement(listKey, key string)
	ZRem(listKey, key string)
	SetString(key, value string)
	Remove(key string)
	GetString(key string) string
	Exist(key string) bool
}
type RedisClient struct {
	client *redis.Client
}

func NewRedisClient(address, password string, db int) *RedisClient {
	rdb := redis.NewClient(&redis.Options{
		Addr:     address,
		Password: password,
		DB:       db,
	})

	return &RedisClient{client: rdb}
}
func (r *RedisClient) ZRangeAdd(listKey, key string) {
	r.client.ZAdd(context.Background(), listKey, &redis.Z{
		Score:  0,
		Member: key,
	})
}
func (r *RedisClient) ZRange(listKey string) []redis.Z {
	slices := r.client.ZRangeWithScores(context.Background(), listKey, 0, -1)
	return slices.Val()
}
func (r *RedisClient) ZRangeScore(listKey, key string) int64 {
	slices := r.ZRange(listKey)
	for _, v := range slices {
		if v.Member == key {
			return int64(v.Score)
		}
	}
	return 0
}
func (r *RedisClient) ZIncrement(listKey, key string) {
	r.client.ZIncrBy(context.Background(), listKey, 1, key)
}

func (r *RedisClient) ZRem(listKey, key string) {
	r.client.ZRem(context.Background(), listKey, key)
}
func (r *RedisClient) SetString(key, value string) {
	r.client.Set(context.Background(), key, value, 0)
}
func (r *RedisClient) Remove(key string) {
	r.client.Del(context.Background(), key)
}
func (r *RedisClient) Exist(key string) bool {
	exist := r.client.Exists(context.Background(), key)
	return exist.Val() > 0
}
func (r *RedisClient) GetString(key string) string {
	strCommand := r.client.Get(context.Background(), key)
	return strCommand.Val()
}
