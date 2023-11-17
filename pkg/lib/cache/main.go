package cache

import (
	"context"
	"errors"
	"fmt"
	"github.com/redis/go-redis/v9"
	"os"
	"time"
)

var cache *redis.Client

func NewCache() error {
	rdb := redis.NewClient(&redis.Options{
		Addr:     fmt.Sprintf("%s:6379", os.Getenv("REDIS_HOST")),
		Password: os.Getenv("REDIS_PASSWORD"),
		DB:       0,
	})

	ctx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	err := rdb.Set(ctx, "test", "", 0).Err()

	if err != nil {
		return errors.New(err.Error())
	}

	delCtx, _ := context.WithDeadline(context.Background(), time.Now().Add(5*time.Second))
	rdb.Del(delCtx, "test")

	cache = rdb

	return nil
}

func Cache() *redis.Client {
	return cache
}
