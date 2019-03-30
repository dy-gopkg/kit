package redis

import (
	"github.com/gomodule/redigo/redis"
	"sync"
	"time"
)

var _redisPool *redis.Pool
var once sync.Once

//单例返回 redisCon Instance
func Do(cmd string, args ...interface{}) (interface{}, error) {
	con := _redisPool.Get()
	defer con.Close()
	return con.Do(cmd, args...)
}

func Init(addr, password string, maxIdle, maxActive int) {
	readTimeout := redis.DialReadTimeout(time.Second * time.Duration(2))
	writeTimeout := redis.DialWriteTimeout(time.Second * time.Duration(2))
	conTimeout := redis.DialConnectTimeout(time.Second * time.Duration(5))
	_redisPool = &redis.Pool{
		MaxIdle:     maxIdle,
		MaxActive:   maxActive,
		IdleTimeout: 0,
		Wait:        true,
		Dial: func() (redis.Conn, error) {
			c, err := redis.Dial("tcp", addr, readTimeout, writeTimeout, conTimeout)
			if err != nil {
				return nil, err
			}
			if len(password) > 0 {
				if _, err := c.Do("AUTH", password); err != nil {
					c.Close()
					return nil, err
				}
			}
			return c, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}
}
