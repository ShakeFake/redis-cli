package utils

import (
	"sync"
)

type Notification struct {
	infos map[string]chan *RedisEvent
	lock  sync.Mutex
}

type NotificationInterface interface {
	NotificationKey(key string) <-chan *RedisEvent
	UnNotificationKey(key string) (error, bool)
}

func (n *Notification) NotificationKey(key string) <-chan *RedisEvent {
	_, ok := n.infos[key]
	// 如果不存在，就创建一个
	if !ok {
		n.infos[key] = make(chan *RedisEvent, 0)
	}
	// 如果存在，直接返回。
	return n.infos[key]
}

func (n *Notification) UnNotificationKey(key string) (error, bool) {
	n.lock.Lock()
	defer n.lock.Unlock()
	_, ok := n.infos[key]
	if ok {
		delete(n.infos, key)
	}
	conn := redisManager.Pool.Get()
	_, err := conn.Do("DEL", key)
	if err != nil {
		return err, false
	}
	return nil, true
}
