package main

import (
	"context"
	"os"
	"time"

	"github.com/mondegor/go-core/mrlog"
	"github.com/mondegor/go-core/mrlog/slog"

	"github.com/mondegor/go-storage/mrredis"
)

func main() {
	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	mrlog.Info(logger, "Create redis connection")

	opts := mrredis.Options{
		Host:         "127.0.0.1",
		Port:         "6379",
		Password:     "123456",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx := context.Background()
	redisAdapter := mrredis.New(logger)

	if err := redisAdapter.Connect(ctx, opts); err != nil {
		mrlog.Fatal(logger, "redisAdapter.Connect()", "error", err)
	}

	defer func() {
		_ = redisAdapter.Close()
	}()

	redisCli, err := redisAdapter.Cli()
	if err != nil {
		mrlog.Fatal(logger, "redisAdapter.Cli()", "error", err)
	}

	key := "my-test-key"
	redisCli.Set(ctx, key, "my-test-value", 1*time.Second)
	value := redisCli.Get(ctx, key).Val()

	mrlog.Info(logger, "value from redis by key", "key", key, "value", value)
}
