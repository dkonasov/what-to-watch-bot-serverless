package session

import (
	"context"
	"os"
	"strconv"

	redis "github.com/go-redis/redis/v8"
)

func Set_step(user string, step int64) error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     os.Getenv("REDIS_URL"),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})
	ctx := context.Background()

	_, err := rdb.HSet(ctx, user, "step", strconv.FormatInt(step, 10)).Result()

	return err
}
