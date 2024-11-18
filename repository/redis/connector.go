package redis

import (
	"context"

	"github.com/redis/go-redis/v9"
)

type Config struct {
	Address  string
	Password string
	DB       int
}

type Client struct {
	appName string
	cli     *redis.Client
}

func New(ctx context.Context, cfg Config, appName string) (*Client, error) {
	rdb := redis.NewClient(&redis.Options{
		Addr:     cfg.Address,
		Password: cfg.Password,
		DB:       cfg.DB,
	})

	if res := rdb.Ping(ctx); res.Err() != nil {
		return nil, res.Err()
	}

	return &Client{cli: rdb, appName: appName}, nil
}
