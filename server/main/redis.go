package main

import (
	"time"

	"github.com/garyburd/redigo/redis"
)

var pool *redis.Pool

func initPool(address string, maxIdle, maxActive int, idleTimeout time.Duration) {

	pool = &redis.Pool{
		MaxIdle:     maxIdle,     //最大空閒鏈接數
		MaxActive:   maxActive,   //表示和數據庫的最大鏈接數
		IdleTimeout: idleTimeout, //最大空閒時間
		Dial: func() (redis.Conn, error) {
			return redis.Dial("tcp", address) //初始化鏈接的代碼 鏈接哪個ip的redis
		},
	}
}
