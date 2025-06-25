package utils

import (
	"encoding/json"
	"fmt"
	"sync"
	"testing"
)

func lineKeys() {
	RedisHost = "r-0xijfeb5nrr4xrqucdpd.redis.rds.aliyuncs.com"
	RedisPort = 6379
	RedisPassword = "JYgKVoaC3Um5IhUN"
}

func testKeys() {
	RedisHost = "192.168.1.1"
	RedisPort = 6379
}

func init() {
	testKeys()
}

type Student struct {
	Name string `json:"name"`
}

func TestRedisManager_Set(t *testing.T) {
	initRedis()
	s := Student{Name: "abc"}
	info, _ := json.Marshal(s)
	err := redisManager.Set("key", (info))
	if err != nil {
		panic(err)
	}
}

func TestRedisManager_Get(t *testing.T) {
	initRedis()
	info, err := redisManager.Get("key")
	if err != nil {
		panic(err)
	}

	var s Student
	json.Unmarshal(info, &s)
	fmt.Println(s)
}

func TestRedisManager_Hset(t *testing.T) {
	initRedis()
	redisManager.Hset()
}

func TestRedisManager_Incr(t *testing.T) {
	initRedis()

	wg := sync.WaitGroup{}

	for i := 0; i < 10000000; i++ {
		wg.Add(1)
		go func() {
			redisManager.Incr()
			wg.Done()
		}()
	}
	wg.Wait()

}

func TestRedisManager_del(t *testing.T) {
	initRedis()
	path := "D:\\myGithub\\redis-cli\\utils\\redis_garbage_key.txt"
	allKeys := ReadRedisKeyFromFile(path)
	for _, key := range allKeys {
		err := redisManager.Del(key)
		if err != nil {
			fmt.Println("delete:", key, "failed:", err)
		}
		fmt.Println("delete:", key, " success")
	}
}
