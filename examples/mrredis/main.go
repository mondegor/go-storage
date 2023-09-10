package main

import (
    "context"
    "time"

    "github.com/mondegor/go-core/mrcore"
    "github.com/mondegor/go-core/mrlog"
    "github.com/mondegor/go-storage/mrredis"
)

func main() {
    logger, _ := mrlog.New("[mrredis]", "debug", false)
    logger.Info("Create redis connection")

    appHelper := mrcore.NewAppHelper(logger)

    opt := mrredis.Options{
        Host: "127.0.0.1",
        Port: "6379",
        Password: "123456",
        ConnTimeout: 10 * time.Second,
    }

    redisClient := mrredis.New()
    err := redisClient.Connect(opt)

    appHelper.ExitOnError(err)
    defer appHelper.Close(redisClient)

    key := "my-test-key"
    redisClient.Cli().Set(context.Background(), key, "my-test-value", 1 * time.Second)
    value := redisClient.Cli().Get(context.Background(), key).Val()

    logger.Info("value from redis by key '%s': %s", key, value)
}
