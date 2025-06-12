package utils

import (
	"errors"
	"fmt"
	"github.com/gomodule/redigo/redis"
	"strings"
	"sync"
	"time"
)

var (
	redisManager     *RedisManager
	notificationList Notification
	// 自定义错误
	notInitRedisManager = errors.New("not init redis manager")
)

func initRedis() {
	redisManager = &RedisManager{}
	open()
}

type RedisManager struct {
	Pool   redis.Pool
	IsOpen bool
	Error  error
	lock   sync.Mutex
}

type operator interface {
	Get(key string) (string, error)
	Set(key string, value string) error
}

// RedisEvent redis keyspace notifications 样例
type RedisEvent struct {
	Number string
	Event  string
	Key    string
}

func GetRedisEvent(number string, event string, key string) *RedisEvent {
	return &(RedisEvent{number, event, key})
}

func open() {
	redisManager.lock.Lock()
	defer redisManager.lock.Unlock()

	redisManager.Pool = redis.Pool{
		Wait:        true,
		MaxActive:   10,
		MaxIdle:     10,
		IdleTimeout: time.Second * 60,
		Dial: func() (redis.Conn, error) {
			if RedisPassword != "" {
				return redis.DialURL(fmt.Sprintf("redis://%v:%v", RedisHost, RedisPort),
					redis.DialPassword(RedisPassword))
			}
			return redis.DialURL(fmt.Sprintf("redis://%v:%v", RedisHost, RedisPort))
		},
	}
	conn := redisManager.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("ping")
	if err != nil {
		redisManager.IsOpen = false
		redisManager.Error = err
	}
}

func GetRedisManager() *RedisManager {
	if redisManager.IsOpen {
		return redisManager
	} else {
		open()
		return redisManager
	}
}

func GetConnection() redis.Conn {
	return redisManager.Pool.Get()
}

func u82b(u8 []uint8) []byte {
	b := make([]byte, len(u8))
	for i, v := range u8 {
		b[i] = byte(v)
	}
	return b
}

func (r *RedisManager) Get(key string) ([]byte, error) {
	conn := r.Pool.Get()
	defer conn.Close()

	result, err := conn.Do("get", key)
	//return u82b(result.([]uint8)), err
	return result.([]byte), err
}

func (r *RedisManager) Set(key string, value interface{}) error {
	conn := r.Pool.Get()
	defer conn.Close()
	_, err := conn.Do("set", key, value)
	return err
}

func (r *RedisManager) Hset() {
	conn := r.Pool.Get()
	defer conn.Close()

	conn.Do("HSET", "k1")

}

func (r *RedisManager) Incr() {
	conn := r.Pool.Get()
	defer conn.Close()

	conn.Do("INCR", "k1")

}

func (r *RedisManager) Del(key string) error {
	conn := r.Pool.Get()
	defer conn.Close()

	_, err := conn.Do("DEL", key)
	return err
}

func ReceiveEvent() error {
	for {
		if redisManager.IsOpen {
			conn := redisManager.Pool.Get()
			psc := redis.PubSubConn{Conn: conn}
			err := psc.PSubscribe("__keyevent*__:*")
			if err != nil {
				return err
			}
		innerLoop:
			for {
				switch message := psc.Receive().(type) {
				case error:
					// 如果是 EOF，进行一次重新获取 conn
					if strings.Contains(message.Error(), "EOF") {
						break innerLoop
					}
					break innerLoop
				case redis.Message:
					subKey := string(message.Data)
					if notification, ok := notificationList.infos[subKey]; ok {
						// __keyevent@0__:set
						subChannel := strings.Split(message.Channel, "@")[1]

						subDatabaseNumber := strings.Split(subChannel, "__:")[0]
						subEvent := strings.Split(subChannel, "__:")[1]
						notification <- GetRedisEvent(subDatabaseNumber, subEvent, subKey)
					}

				case redis.Subscription:
					fmt.Printf("Subscription: kind is %s, channel is %s, count is %d", message.Kind, message.Channel, message.Count)
				}
			}
		} else {
			break
		}
	}
	return notInitRedisManager
}
