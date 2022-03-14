package redis

import (
	"fmt"

	"github.com/go-redis/redis/v8"
	"github.com/stretchr/testify/mock"
)

type RedisClientMock struct {
	ZRangeData []redis.Z
	mock.Mock
}

func (r *RedisClientMock) ZRangeAdd(listKey, key string) {
	fmt.Println("ZRangeAdd called")

}
func (r *RedisClientMock) ZRange(listKey string) []redis.Z {

	return []redis.Z{}
}
func (r *RedisClientMock) ZRangeScore(listKey, key string) int64 {

	return 0
}
func (r *RedisClientMock) ZIncrement(listKey, key string) {

}

func (r *RedisClientMock) ZRem(listKey, key string) {

}
func (r *RedisClientMock) SetString(key, value string) {
	fmt.Println("SetString called")
}
func (r *RedisClientMock) Remove(key string) {
}
func (r *RedisClientMock) Exist(key string) bool {
	if key == "exist" {
		return true
	}
	return false
}
func (r *RedisClientMock) GetString(key string) string {
	return ""
}
