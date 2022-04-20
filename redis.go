package main

import (
	"fmt"
	"github.com/UncleDeron/frogChat-server/utils"
	"github.com/aceld/zinx/zlog"
	"time"

	"github.com/garyburd/redigo/redis"
)

func initPool(rc *utils.RedisConfig) *redis.Pool {
	pool := &redis.Pool{
		Dial: func() (redis.Conn, error) { // 初始化连接函数
			c, err := redis.Dial("tcp", fmt.Sprintf("%s:%d", rc.Host, rc.Port))
			if err != nil {
				zlog.Error(err)
				return nil, err
			}
			if rc.Password != "" {
				if _, err := c.Do("AUTH", rc.Password); err != nil {
					c.Close()
					zlog.Error(err)
					return nil, err
				}
			}
			if _, err := c.Do("SELECT", rc.DB); err != nil {
				c.Close()
				zlog.Error(err)
				return nil, err
			}
			return c, nil
		},
		MaxIdle:     rc.MaxIdle,                                  // 最大空闲连接数
		MaxActive:   rc.MaxActive,                                // 和数据库的最大连接数, 0 表示不限制
		IdleTimeout: time.Duration(rc.IdleTimeout) * time.Second, // 最大空闲时间, 超过该时间, 空闲的连接将被关闭, 0 表示不关闭连接
	}
	return pool
}
