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
	connectionName = "Minio"
)

type (
	// ConnAdapter - адаптер для работы с Minio клиентом.
	ConnAdapter struct {
		conn          *minio.Client
		tracer        mrtrace.Tracer
		createBuckets bool // if not exists
		mimeTypes     *mime.TypeList
	}

	// Options - опции для создания соединения для ConnAdapter.
	Options struct {
		DSN      string // если указано, то Host, Port не используются
		Host     string
		Port     string
		UseSSL   bool
		User     string
		Password string
	}
)

var errSystemStorageConnectionIsBusy = errors.NewSystemProto("connection is busy")

// New - создаёт объект ConnAdapter.
func New(createBuckets bool, mimeTypes *mime.TypeList, tracer mrtrace.Tracer) *ConnAdapter {
	return &ConnAdapter{
		createBuckets: createBuckets,
		mimeTypes:     mimeTypes,
		tracer:        tracer,
	}
}

// Connect - создаёт соединение с указанными опциями.
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

// Ping - сообщает, установлено ли соединение и является ли оно стабильным.
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

// InitBucket - инициализирует бакет: проверяет что он существует, и если нет,
// то или создаёт его (если разрешено) или выдаёт ошибку.
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

// Cli - возвращается нативный объект, с которым работает данный адаптер.
func (c *ConnAdapter) Cli() (*minio.Client, error) {
	if c.conn == nil {
		return nil, errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	return c.conn, nil
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return errors.ErrInternalStorageConnectionIsNotOpened.New("source", connectionName)
	}

	c.conn = nil

	return nil
}
