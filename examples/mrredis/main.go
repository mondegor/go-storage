package main

import (
	"context"
	"os"
	"time"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrtrace/logtracer"

	"github.com/mondegor/go-storage/mrredis"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

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
	redisAdapter := mrredis.New(logtracer.New(logger))

	if err := redisAdapter.Connect(ctx, opts); err != nil {
		mrlog.Fatal(logger, "redisAdapter.Connect()", "error", err) // mrlog.Fatal
	}

	defer redisAdapter.Close()

	redisCli, err := redisAdapter.Cli()
	if err != nil {
		mrlog.Fatal(logger, "redisAdapter.Cli()", "error", err) // mrlog.Fatal
	}

	key := "my-test-key"
	redisCli.Set(ctx, key, "my-test-value", 1*time.Second)
	value := redisCli.Get(ctx, key).Val()

	mrlog.Info(logger, "value from redis by key", "key", key, "value", value)
}
