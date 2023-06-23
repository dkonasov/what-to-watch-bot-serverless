package session

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func Set_active_list(user string, active_list uint) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.Background()

	_, err := rdb.HSet(ctx, user, "active_list", strconv.FormatUint(uint64(active_list), 10)).Result()

	return err
}
