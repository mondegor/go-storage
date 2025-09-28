package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrrabbitmq"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	logger.Info("Create rabbitmq connection")

	opts := mrrabbitmq.Options{
		Host:     "127.0.0.1",
		Port:     "5672",
		User:     "admin",
		Password: "123456",
	}

	ctx := context.Background()
	rabbitAdapter := mrrabbitmq.New()

	if err := rabbitAdapter.Connect(ctx, opts); err != nil {
		logger.Error("rabbitAdapter.Connect()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	defer rabbitAdapter.Close()

	logger.Info("Create rabbitmq channel")

	rabbitCli, err := rabbitAdapter.Cli()
	if err != nil {
		logger.Error("rabbitAdapter.Cli()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	rabbitChannel, err := rabbitCli.Channel()
	if err != nil {
		logger.Error("rabbitCli.Channel()", "error", err) // mrlog.Fatal
		os.Exit(1)
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
		logger.Error("rabbitChannel.QueueDeclare()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}
}
