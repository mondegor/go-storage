package mrrabbitmq

import (
    "fmt"

    "github.com/mondegor/go-storage/mrstorage"

    amqp "github.com/rabbitmq/amqp091-go"
)

// go get github.com/rabbitmq/amqp091-go@v1.8.1

const connectionName = "rabbitmq"

type (
    Connection struct {
        conn *amqp.Connection
    }

    Options struct {
        Host string
        Port string
        User string
        Password string
    }
)

func New() *Connection {
    return &Connection{}
}

func (c *Connection) Connect(opt Options) error {
    if c.conn != nil {
        return mrstorage.ErrFactoryConnectionIsAlreadyCreated.New(connectionName)
    }

    conn, err := amqp.Dial(fmt.Sprintf("amqp://%s:%s@%s:%s/", opt.User, opt.Password, opt.Host, opt.Port))

    if err != nil {
        return mrstorage.ErrFactoryConnectionFailed.Wrap(err, connectionName)
    }

    c.conn = conn

    return nil
}

func (c *Connection) Cli() *amqp.Connection {
    return c.conn
}

func (c *Connection) Close() error {
    if c.conn == nil {
        return mrstorage.ErrFactoryConnectionIsNotOpened.New(connectionName)
    }

    err := c.conn.Close()

    if err != nil {
        return mrstorage.ErrFactoryConnectionFailed.Wrap(err, connectionName)
    }

    c.conn = nil

    return nil
}
