package mrminio

import (
	"context"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/mime"
)

// go get -u github.com/minio/minio-go/v7
// https://min.io/docs/minio/linux/developers/go/API.html

const (
	// connectionName - имя подключения для логирования и трассировки.
	connectionName = "Minio"
)

type (
	// ConnAdapter - адаптер для работы с MinIO клиентом (S3-совместимое хранилище).
	// Предоставляет методы для подключения, проверки работоспособности,
	// инициализации бакетов и получения нативного клиента MinIO.
	ConnAdapter struct {
		conn          *minio.Client  // conn - нативный клиент MinIO
		tracer        mrtrace.Tracer // tracer - трассировщик для логирования операций
		createBuckets bool           // createBuckets - флаг автоматического создания бакетов, если они не существуют
		mimeTypes     *mime.TypeList // mimeTypes - список MIME-типов для определения типа контента
	}

	// Options - опции для создания соединения в ConnAdapter.
	// Позволяет подключаться либо по DSN, либо по отдельным параметрам Host/Port.
	Options struct {
		DSN      string // DSN - строка подключения в формате "host:port" (если указана, Host и Port игнорируются)
		Host     string // Host - адрес сервера MinIO (используется, если DSN не указан)
		Port     string // Port - порт сервера MinIO (используется, если DSN не указан)
		UseSSL   bool   // UseSSL - флаг использования SSL/TLS соединения
		User     string // User - имя пользователя для аутентификации
		Password string // Password - пароль для аутентификации
	}
)

var errSystemStorageConnectionIsBusy = errors.NewSystemProto("connection is busy")

// New - создаёт объект ConnAdapter без активного соединения.
// Параметры:
//   - createBuckets - если true, автоматически создавать бакеты при их отсутствии;
//   - mimeTypes - список MIME-типов для определения типа контента;
//   - tracer - трассировщик для логирования операций.
func New(createBuckets bool, mimeTypes *mime.TypeList, tracer mrtrace.Tracer) *ConnAdapter {
	return &ConnAdapter{
		createBuckets: createBuckets,
		mimeTypes:     mimeTypes,
		tracer:        tracer,
	}
}

// Connect - устанавливает соединение с сервером MinIO по указанным опциям.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return errors.ErrInternalStorageConnectionIsAlreadyCreated.New("source", connectionName)
	}

	if opts.DSN == "" {
		opts.DSN = fmt.Sprintf("%s:%s", opts.Host, opts.Port)
	}

	conn, err := minio.New(
		opts.DSN,
		&minio.Options{
			Creds:  credentials.NewStaticV4(opts.User, opts.Password, ""),
			Secure: opts.UseSSL,
		},
	)
	if err != nil {
		return errors.ErrSystemStorageConnectionFailed.Wrap(err, "source", connectionName)
	}

	c.conn = conn

	return nil
}

// Ping - проверяет работоспособность соединения с MinIO.
func (c *ConnAdapter) Ping(_ context.Context) error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	// TODO: желательно вызывать cancel(), если внешний контекст отменился, без этого HealthCheck вылетит по таймауту только через 3 секунды

	cancel, err := c.conn.HealthCheck(time.Hour) // 1 час - нужно чтобы внутренняя горутина гарантирована не вызвалась
	if err != nil {
		return errSystemStorageConnectionIsBusy.Wrap(err, "source", connectionName)
	}
	defer cancel()

	if c.conn.IsOffline() {
		return errors.ErrSystemStorageConnectionFailed.WithDetails(
			"expected status 'online' but found 'offline'",
			"source", connectionName,
		)
	}

	return nil
}

// InitBucket - инициализирует бакет: проверяет его существование и при необходимости создаёт.
// Возвращает true, если бакет был создан, false - если уже существовал.
// Если бакет не существует и createBuckets=false, возвращает ошибку.
func (c *ConnAdapter) InitBucket(ctx context.Context, bucketName string) (bool, error) {
	exists, err := c.conn.BucketExists(ctx, bucketName)
	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	if !c.createBuckets {
		return false, fmt.Errorf("bucket not exists (name='%s')", bucketName)
	}

	err = c.conn.MakeBucket(
		ctx,
		bucketName,
		minio.MakeBucketOptions{}, // "ru-central1"
	)
	if err != nil {
		return false, err
	}

	return true, nil
}

// Cli - возвращает нативный клиент MinIO для прямого доступа к API.
func (c *ConnAdapter) Cli() (*minio.Client, error) {
	if c.conn == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.conn, nil
}

// Close - закрывает соединение с MinIO и очищает ссылку на клиент.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	c.conn = nil

	return nil
}
