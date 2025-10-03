package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrrabbitmq"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	mrlog.Info(logger, "Create rabbitmq connection")

	opts := mrrabbitmq.Options{
		Host:     "127.0.0.1",
		Port:     "5672",
		User:     "admin",
		Password: "123456",
	}

	ctx := context.Background()
	rabbitAdapter := mrrabbitmq.New()

	if err := rabbitAdapter.Connect(ctx, opts); err != nil {
		mrlog.Fatal(logger, "rabbitAdapter.Connect()", "error", err) // mrlog.Fatal
	}

	defer rabbitAdapter.Close()

	mrlog.Info(logger, "Create rabbitmq channel")

	rabbitCli, err := rabbitAdapter.Cli()
	if err != nil {
		mrlog.Fatal(logger, "rabbitAdapter.Cli()", "error", err) // mrlog.Fatal
	}

	rabbitChannel, err := rabbitCli.Channel()
	if err != nil {
		mrlog.Fatal(logger, "rabbitCli.Channel()", "error", err) // mrlog.Fatal
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
		mrlog.Fatal(logger, "rabbitChannel.QueueDeclare()", "error", err) // mrlog.Fatal
	}
}
