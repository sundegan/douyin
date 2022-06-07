package dao

import (
	"fmt"
	"github.com/go-redis/redis/v8"
)

// Del 根据key删除键值对(可删除多对)
func Del(rdb *redis.Client, keys []string) int64 {
	val, err := rdb.Del(Ctx, keys...).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// Exists 检查某个key是否存在
func Exists(rdb *redis.Client, key string) bool {
	n, err := rdb.Exists(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	if n > 0 {
		return true
	} else {
		return false
	}
}

// DBSize 查看当前数据库键值对数量
func DBSize(rdb *redis.Client) int64 {
	size, err := rdb.DBSize(Ctx).Result()
	if err != nil {
		panic(err)
	}
	return size
}

// FlushDB 清空当前数据库
func FlushDB(rdb *redis.Client) {
	res, err := rdb.FlushDB(Ctx).Result()
	if err != nil {
		panic(err)
	}
	fmt.Println(res) // OK
}

// HSet 给<key>集合中的<field>字段赋值<value>
func HSet(rdb *redis.Client, key, field, name string) {
	rdb.HSet(Ctx, key, field, name)
}

// HGet 从<key>集合<field>字段取出value
func HGet(rdb *redis.Client, key, field string) string {
	val, err := rdb.HGet(Ctx, key, field).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HExists 判断哈希表key对应的field是否存在
func HExists(rdb *redis.Client, key, field string) bool {
	val, err := rdb.HExists(Ctx, key, field).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HKeys 列出该key的所有field字段
func HKeys(rdb *redis.Client, key string) []string {
	val, err := rdb.HKeys(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HVals 列出该key的所有value值
func HVals(rdb *redis.Client, key string) []string {
	val, err := rdb.HVals(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HGetAll 返回哈希表key对应的所有字段和值
func HGetAll(rdb *redis.Client, key string) map[string]string {
	val, err := rdb.HGetAll(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HDel 删除哈希表key中的一个或多个指定字段，不存在的字段将被忽略
func HDel(rdb *redis.Client, key, field string) int64 {
	val, err := rdb.HDel(Ctx, key, field).Result()
	if err != nil {
		panic(err)
	}
	return val
}

// HLen 获取哈希表中key中字段的数量
func HLen(rdb *redis.Client, key string) int64 {
	length, err := rdb.HLen(Ctx, key).Result()
	if err != nil {
		panic(err)
	}
	return length
}
