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
	// ConnAdapter - адаптер для работы с Rabbitmq клиентом.
	ConnAdapter struct {
		conn *amqp.Connection
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		Host     string
		Port     string
		User     string
		Password string
	}
)

// New - создаёт объект ConnAdapter.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - создаёт соединение с указанными опциями.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return mrcore.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
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
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = conn

	return nil
}

// Cli - возвращается нативный объект, с которым работает данный адаптер.
func (c *ConnAdapter) Cli() *amqp.Connection {
	return c.conn
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = nil

	return nil
}
