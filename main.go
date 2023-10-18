package main

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
	"wilikidi/redis-cli/utils"
	"wilikidi/redis-cli/utils/random"
)

type Sleep interface {
	Sleep(time time.Duration)
}

type First struct {
	sleepTime time.Duration
}

func (f *First) Sleep(time time.Duration) {
	f.sleepTime = time
}

type Second struct {
	sleepTime time.Duration
	sleep     func(time time.Duration)
}

func (s *Second) Sleep(time time.Duration) {
	fmt.Println("休眠5s")
}

func main() {
	totalTimes := 99999

	signal := make(chan int, 80)

	successTimes := 0
	var mu sync.Mutex

	for i := 0; i < totalTimes; i++ {
		signal <- 1
		go func() {
			coon := utils.GetConnection()
			defer func() {
				// 先关
				coon.Close()
				<-signal
			}()
			ruuid := random.GetRandomUUID()
			reply, err := coon.Do("set", ruuid, ruuid+"value")
			// 10 个就会造成资源耗尽。如果需要处理，则需要开多进程。控制 获取资源量。
			if err != nil {
				panic(err)
			}
			if k, err := redis.String(reply, err); err != nil && k == "OK" {
				// 这块有并发问题
				mu.Lock()
				successTimes += 1
				mu.Unlock()
			}
		}()

	}
	fmt.Printf("all times %v, success times %v", totalTimes, successTimes)
}
