package db

import (
	"fmt"
	"github.com/go-redis/redis/v8"
	"sync"
	"time"
)

var RedisClientMap = map[string]*redis.Client{}
var redisLock sync.Mutex

type RedisOption struct {
	Address  string
	Password string
	Db       int
}

// GetRedisInstance 输入一个redis的配置选项,返回一个redis的客户端
func GetRedisInstance(redisOpt RedisOption) *redis.Client {
	address := redisOpt.Address
	db := redisOpt.Db
	password := redisOpt.Password
	addr := fmt.Sprintf("%s", address)
	redisLock.Lock()
	if redisCli, ok := RedisClientMap[addr]; ok { //是否已存在相同ip的redis连接
		return redisCli
	}
	client := redis.NewClient(&redis.Options{
		Addr:       addr,
		Password:   password,
		DB:         db,
		MaxConnAge: 20 * time.Second,
	})
	
	RedisClientMap[addr] = client //存入map
	redisLock.Unlock()
	return RedisClientMap[addr]
}
