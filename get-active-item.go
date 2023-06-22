package main

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func get_active_item(user string) (uint, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.Background()

	raw_value, err := rdb.HGet(ctx, user, "active_item").Result()

	if err == nil || err.Error() == "redis: nil" {
		if err != nil {
			return 0, nil
		}

		item, err := strconv.ParseUint(raw_value, 10, 32)

		if err == nil {
			return uint(item), nil
		}
	}

	return 0, err
}
