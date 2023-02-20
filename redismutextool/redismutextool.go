package redismutextool

import (
	"time"

	"github.com/go-redsync/redsync/v4"
	"github.com/go-redsync/redsync/v4/redis/goredis/v9"
	goredislib "github.com/redis/go-redis/v9"
)

type RedisMutexTool struct {
	Client goredislib.UniversalClient
	Mutex  *redsync.Mutex
}

func NewRedisMutexTool(addrs []string, mutexName string, expire time.Duration) *RedisMutexTool {
	var redisMutexTool RedisMutexTool
	client := goredislib.NewClient(&goredislib.Options{
		Addr: addrs[0],
	})
	redisMutexTool.Client = client
	if len(addrs) >= 2 {
		clusterClient := goredislib.NewClusterClient(&goredislib.ClusterOptions{
			Addrs: addrs,
		})
		redisMutexTool.Client = clusterClient
	}
	pool := goredis.NewPool(redisMutexTool.Client)
	rs := redsync.New(pool)
	redisMutexTool.Mutex = rs.NewMutex(mutexName, redsync.WithExpiry(expire))
	return &redisMutexTool
}

func (rm *RedisMutexTool) Lock() error {
	return rm.Mutex.Lock()
}

func (rm *RedisMutexTool) UnLock() (bool, error) {
	return rm.Mutex.Unlock()
}

// Extend 续约
func (rm *RedisMutexTool) Extend() (bool, error) {
	return rm.Mutex.Extend()
}

func (rm *RedisMutexTool) Close() {
	rm.Client.Close()
}
