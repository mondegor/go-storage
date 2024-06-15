package main

import (
	"context"

	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrlog/mrlogbase"

	"github.com/mondegor/go-storage/mrrabbitmq"
)

func main() {
	logger := mrlogbase.New(mrlog.TraceLevel)
	ctx := mrlog.WithContext(context.Background(), logger)

	logger.Info().Msg("Create rabbitmq connection")

	opts := mrrabbitmq.Options{
		Host:     "127.0.0.1",
		Port:     "5672",
		User:     "admin",
		Password: "123456",
	}

	rabbitAdapter := mrrabbitmq.New()

	if err := rabbitAdapter.Connect(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("rabbitAdapter.Connect() error")
	}

	defer rabbitAdapter.Close()

	logger.Info().Msg("Create rabbitmq channel")

	rabbitChannel, err := rabbitAdapter.Cli().Channel()
	if err != nil {
		logger.Fatal().Err(err).Msg("rabbitAdapter.Cli().Channel() error")
	}

	defer rabbitChannel.Close()

	_, err = rabbitChannel.QueueDeclare(
		"my.test.queue", // name
		true,            // durable
		false,           // autoDelete
		true,            // exclusive
		false,           // noWait
		nil,             // args
	)
	if err != nil {
		logger.Fatal().Err(err).Msg("rabbitChannel.QueueDeclare() error")
	}
}
