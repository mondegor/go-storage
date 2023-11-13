package main

import (
    "context"
    "time"

    "github.com/mondegor/go-storage/mrredis"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

func main() {
    logger := mrcore.Log().With("mrredis")
    logger.Info("Create redis connection")

    appHelper := mrtool.NewAppHelper(logger)

    opt := mrredis.Options{
        Host: "127.0.0.1",
        Port: "6379",
        Password: "123456",
        ConnTimeout: 10 * time.Second,
    }

    redisAdapter := mrredis.New()
    err := redisAdapter.Connect(opt)

    appHelper.ExitOnError(err)
    defer appHelper.Close(redisAdapter)

    key := "my-test-key"
    redisAdapter.Cli().Set(context.Background(), key, "my-test-value", 1 * time.Second)
    value := redisAdapter.Cli().Get(context.Background(), key).Val()

    logger.Info("value from redis by key '%s': %s", key, value)
}
