package mrminio

import (
	"context"
	"errors"
	"fmt"
	"time"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrtrace"
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
		mimeTypes     *extfile.MimeTypeList
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

// New - создаёт объект ConnAdapter.
func New(createBuckets bool, mimeTypes *extfile.MimeTypeList, tracer mrtrace.Tracer) *ConnAdapter {
	return &ConnAdapter{
		createBuckets: createBuckets,
		mimeTypes:     mimeTypes,
		tracer:        tracer,
	}
}

// Connect - создаёт соединение с указанными опциями.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return mr.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
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
		return mr.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = conn

	return nil
}

// Ping - сообщает, установлено ли соединение и является ли оно стабильным.
func (c *ConnAdapter) Ping(_ context.Context) error {
	if c.conn == nil {
		return mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	// TODO: желательно вызывать cancel(), если внешний контекст отменился, без этого HealthCheck вылетит по таймауту только через 3 секунды

	cancel, err := c.conn.HealthCheck(time.Hour) // 1 час - нужно чтобы внутренняя горутина гарантирована не вызвалась
	if err != nil {
		return mr.ErrStorageConnectionIsBusy.Wrap(err, connectionName)
	}
	defer cancel()

	if c.conn.IsOffline() {
		return mr.ErrStorageConnectionFailed.Wrap(errors.New("expected status 'online' but found 'offline'"), connectionName)
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
		return false, fmt.Errorf("bucket with name '%s' not exists", bucketName)
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
		return nil, mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	return c.conn, nil
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mr.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.conn = nil

	return nil
}
