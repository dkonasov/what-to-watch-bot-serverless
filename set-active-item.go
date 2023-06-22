package main

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func set_active_item(user string, active_item uint) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.Background()

	_, err := rdb.HSet(ctx, user, "active_item", strconv.FormatUint(uint64(active_item), 10)).Result()

	return err
}
