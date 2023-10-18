package utils

import "testing"

func TestRedisManager_Set(t *testing.T) {
	redisManager.Set("key", "value")
}
