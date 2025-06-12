package utils

import (
	"fmt"
	"testing"
)

// todo: 此处需要换成mock
func TestReadRedisKeyFromFile(t *testing.T) {
	path := "D:\\myGithub\\redis-cli\\utils\\redis_garbage_key.txt"
	keys := ReadRedisKeyFromFile(path)
	for _, key := range keys {
		fmt.Println(key)
	}
}
