package main

import (
	"context"
	"os"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrredis"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	logger.Info("Create redis connection")

	opts := mrredis.Options{
		Host:         "127.0.0.1",
		Port:         "6379",
		Password:     "123456",
		ReadTimeout:  10 * time.Second,
		WriteTimeout: 10 * time.Second,
	}

	ctx := context.Background()
	redisAdapter := mrredis.New(mrlog.NewDebugTracer(l))

	if err := redisAdapter.Connect(ctx, opts); err != nil {
		logger.Error("redisAdapter.Connect()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	defer redisAdapter.Close()

	redisCli, err := redisAdapter.Cli()
	if err != nil {
		logger.Error("redisAdapter.Cli()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	key := "my-test-key"
	redisCli.Set(ctx, key, "my-test-value", 1*time.Second)
	value := redisCli.Get(ctx, key).Val()

	logger.Info("value from redis by key", "key", key, "value", value)
}
