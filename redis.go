package main

import "github.com/go-redis/redis"

var redisdb redis.Client

type Msg struct {
	Status   	string
	Stantion	string
	Text 		string
}

func initRedis() redis.Client{
	redisdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Redis.Addr,
		Password: cfg.Redis.Password,
		DB:       cfg.Redis.Db,
	})

	return *redisdb
}
