package mrrabbitmq

import (
	"context"
	"fmt"

	"github.com/mondegor/go-webcore/mrcore"
	amqp "github.com/rabbitmq/amqp091-go"
)

// go get github.com/rabbitmq/amqp091-go@v1.8.1

const (
	connectionName = "Rabbitmq"
)

type (
	ConnAdapter struct {
		conn *amqp.Connection
	}

	Options struct {
		Host     string
		Port     string
		User     string
		Password string
	}
)

func New() *ConnAdapter {
	return &ConnAdapter{}
}

func (c *ConnAdapter) Connect(ctx context.Context, opts Options) error {
	if c.conn != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	conn, err := amqp.Dial(
		fmt.Sprintf(
			"amqp://%s:%s@%s:%s/",
			opts.User,
			opts.Password,
			opts.Host,
			opts.Port,
		),
	)

	if err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = conn

	return nil
}

func (c *ConnAdapter) Cli() *amqp.Connection {
	return c.conn
}

func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = nil

	return nil
}
