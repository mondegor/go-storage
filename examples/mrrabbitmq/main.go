package main

import (
	"github.com/mondegor/go-storage/mrrabbitmq"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

func main() {
	logger := mrcore.DefaultLogger().With("mrrabbitmq")
	mrcore.SetDefaultLogger(logger)

	logger.Info("Create rabbitmq connection")

	appHelper := mrtool.NewAppHelper(logger)

	opt := mrrabbitmq.Options{
		Host:     "127.0.0.1",
		Port:     "5672",
		User:     "admin",
		Password: "123456",
	}

	rabbitAdapter := mrrabbitmq.New()
	err := rabbitAdapter.Connect(opt)

	appHelper.ExitOnError(err)
	defer appHelper.Close(rabbitAdapter)

	logger.Info("Create rabbitmq channel")

	rabbitChannel, err := rabbitAdapter.Cli().Channel()
	appHelper.ExitOnError(err)

	_, err = rabbitChannel.QueueDeclare(
		"my.test.queue", // name
		true,            // durable
		false,           // autoDelete
		true,            // exclusive
		false,           // noWait
		nil,             // args
	)

	appHelper.ExitOnError(err)
}
