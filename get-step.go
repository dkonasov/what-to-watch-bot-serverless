package main

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func get_step(user string) (int64, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.Background()

	raw_step, err := rdb.HGet(ctx, user, "step").Result()

	if err == nil || err.Error() == "redis: nil" {
		if err != nil {
			return 0, nil
		}

		step, err := strconv.ParseInt(raw_step, 10, 64)

		if err == nil {
			return step, nil
		}
	}

	return 0, err
}
