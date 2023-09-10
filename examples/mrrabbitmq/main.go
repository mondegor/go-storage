package main

import (
    "github.com/mondegor/go-core/mrcore"
    "github.com/mondegor/go-core/mrlog"
    "github.com/mondegor/go-storage/mrrabbitmq"
)

func main() {
    logger, _ := mrlog.New("[mrrabbitmq]", "debug", false)
    logger.Info("Create rabbitmq connection")

    appHelper := mrcore.NewAppHelper(logger)

    opt := mrrabbitmq.Options{
        Host: "127.0.0.1",
        Port: "5672",
        User: "admin",
        Password: "123456",
    }

    rabbitClient := mrrabbitmq.New()
    err := rabbitClient.Connect(opt)

    appHelper.ExitOnError(err)
    defer appHelper.Close(rabbitClient)

    logger.Info("Create rabbitmq channel")

    rabbitChannel, err := rabbitClient.Cli().Channel()
    appHelper.ExitOnError(err)

    _, err = rabbitChannel.QueueDeclare(
        "my.test.queue", // name
        true, // durable
        false, // autoDelete
        true, // exclusive
        false, // noWait
        nil, // args
    )

    appHelper.ExitOnError(err)
}
