package database

import (
	"fmt"
	"github.com/gomodule/redigo/redis"
	"review-bot/config"
	"time"
)

var RedisPool *redis.Pool

func RedisSetup() {
	RedisPool = &redis.Pool{
		MaxIdle:     10,
		MaxActive:   100,
		IdleTimeout: 240 * time.Second,
		Dial: func() (redis.Conn, error) {
			conn, err := redis.Dial("tcp", config.RedisConfig.Host,
				redis.DialPassword(config.RedisConfig.Password),
				redis.DialDatabase(config.RedisConfig.Db))
			if err != nil {
				panic(err.Error())
			}
			return conn, err
		},
		TestOnBorrow: func(c redis.Conn, t time.Time) error {
			_, err := c.Do("PING")
			return err
		},
	}

	fmt.Println("Redis connection pool established")
}
