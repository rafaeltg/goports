package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

// Database represents a redis Database.
type Database struct {
	redis.Client
}

func NewDatabase(ctx context.Context, addr string) (*Database, error) {
	client := redis.NewClient(&redis.Options{
		Addr:     addr,
		Password: "",
		DB:       0,
	})

	_, err := client.Ping(ctx).Result()
	if err != nil {
		return nil, err
	}

	return &Database{*client}, nil
}
