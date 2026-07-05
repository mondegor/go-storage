package mrrabbitmq

import (
	"context"
	"fmt"

	"github.com/mondegor/go-core/errors"
	amqp "github.com/rabbitmq/amqp091-go"
)

// go get github.com/rabbitmq/amqp091-go@v1.8.1

const (
	// connectionName - имя подключения для логирования и трассировки.
	connectionName = "Rabbitmq"
)

type (
	// ConnAdapter - адаптер для работы с RabbitMQ клиентом (AMQP 0.9.1).
	// Предоставляет методы для подключения, получения нативного соединения и закрытия.
	ConnAdapter struct {
		conn *amqp.Connection
	}

	// Options - опции для создания соединения в ConnAdapter.
	// Подключение формируется по схеме: amqp://User:Password@Host:Port/
	Options struct {
		Host     string // Host - адрес сервера RabbitMQ
		Port     string // Port - порт сервера RabbitMQ (обычно 5672)
		User     string // User - имя пользователя для аутентификации
		Password string // Password - пароль для аутентификации
	}
)

// New - создаёт объект ConnAdapter без активного соединения.
func New() *ConnAdapter {
	return &ConnAdapter{}
}

// Connect - устанавливает соединение с сервером RabbitMQ по указанным опциям.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return errors.ErrInternalStorageConnectionIsAlreadyCreated.New("source", connectionName)
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
		return errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	c.conn = conn

	return nil
}

// Cli - возвращает нативное соединение amqp.Connection для прямого доступа к API RabbitMQ.
// Позволяет создавать каналы (Channel) для публикации и потребления сообщений.
func (c *ConnAdapter) Cli() (*amqp.Connection, error) {
	if c.conn == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.conn, nil
}

// Close - закрывает соединение с RabbitMQ и очищает ссылку на соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	if err := c.conn.Close(); err != nil {
		return errors.ErrSystemStorageFailedToClose.Wrap(err, "source", connectionName)
	}

	c.conn = nil

	return nil
}
