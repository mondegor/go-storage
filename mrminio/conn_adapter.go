package mrminio

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
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
		createBuckets bool // if not exists
		mimeTypes     *mrlib.MimeTypeList
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
func New(createBuckets bool, mimeTypes *mrlib.MimeTypeList) *ConnAdapter {
	return &ConnAdapter{
		createBuckets: createBuckets,
		mimeTypes:     mimeTypes,
	}
}

// Connect - создаёт соединение с указанными опциями.
func (c *ConnAdapter) Connect(_ context.Context, opts Options) error {
	if c.conn != nil {
		return mrcore.ErrStorageConnectionIsAlreadyCreated.New(connectionName)
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
		return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = conn

	return nil
}

// Ping - проверяет работоспособность соединения.
func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.conn == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	// :TODO: найти способ пинга лучше
	if _, err := c.conn.BucketExists(ctx, "bucket-ping"); err != nil {
		if strings.Contains(err.Error(), "connection") {
			return mrcore.ErrStorageConnectionFailed.Wrap(err, connectionName)
		}
	}

	return nil
}

// InitBucket - comment method.
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
func (c *ConnAdapter) Cli() *minio.Client {
	return c.conn
}

// Close - закрывает текущее соединение.
func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.ErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.conn = nil

	return nil
}
