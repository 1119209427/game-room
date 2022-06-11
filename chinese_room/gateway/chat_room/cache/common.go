package cache

import (
	"github.com/go-redis/redis"
	"time"
)
var Redis *redis.Client

// InitClient 初始化连接
func InitClient() (err error) {
	Redis = redis.NewClient(&redis.Options{
		Addr:     "localhost:6379",
		Password: "", // no password set
		DB:       0,  // use default DB
	})

	_, err = Redis.Ping().Result()
	if err != nil {
		return err
	}
	return nil

}

func Set(key,value string,exp time.Duration)error{
	 err:=Redis.Set(key,value,exp).Err()
	 return err

}

