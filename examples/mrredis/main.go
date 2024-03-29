package main

import (
	"context"
	"time"

	"github.com/mondegor/go-storage/mrredis"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	ctx := mrlog.WithContext(context.Background(), logger)

	logger.Info().Msg("Create redis connection")

	opts := mrredis.Options{
		Host:        "127.0.0.1",
		Port:        "6379",
		Password:    "123456",
		ConnTimeout: 10 * time.Second,
	}

	redisAdapter := mrredis.New()

	if err := redisAdapter.Connect(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("redisAdapter.Connect() error")
	}

	defer redisAdapter.Close()

	key := "my-test-key"
	redisAdapter.Cli().Set(ctx, key, "my-test-value", 1*time.Second)
	value := redisAdapter.Cli().Get(ctx, key).Val()

	logger.Info().Msgf("value from redis by key '%s': %s", key, value)
}
