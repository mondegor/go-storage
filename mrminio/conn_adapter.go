package mrminio

import (
	"context"
	"fmt"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/minio/minio-go/v7/pkg/credentials"
	"github.com/mondegor/go-webcore/mrcore"
)

// go get -u github.com/minio/minio-go/v7
// https://min.io/docs/minio/linux/developers/go/API.html

const (
	connectionName = "minio"
)

type (
	ConnAdapter struct {
		conn *minio.Client
	}

	Options struct {
		Host     string
		Port     string
		UseSSL   bool
		User     string
		Password string
	}
)

func New() *ConnAdapter {
	return &ConnAdapter{}
}

func (c *ConnAdapter) Connect(opt Options) error {
	if c.conn != nil {
		return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
	}

	conn, err := minio.New(
		getURL(&opt),
		&minio.Options{
			Creds:  credentials.NewStaticV4(opt.User, opt.Password, ""),
			Secure: opt.UseSSL,
		},
	)

	if err != nil {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	c.conn = conn

	return nil
}

func (c *ConnAdapter) Ping(ctx context.Context) error {
	if c.conn == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	_, err := c.conn.BucketExists(ctx, "bucket-ping") // :TODO: возможно есть способ пинга лучше

	if err != nil && strings.Contains(err.Error(), "connection") {
		return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
	}

	return nil
}

func (c *ConnAdapter) InitBucket(ctx context.Context, bucketName string, createIfNotExists bool) (bool, error) {
	exists, err := c.conn.BucketExists(ctx, bucketName)

	if err != nil {
		return false, err
	}

	if exists {
		return false, nil
	}

	if !createIfNotExists {
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

func (c *ConnAdapter) Cli() *minio.Client {
	return c.conn
}

func (c *ConnAdapter) Close() error {
	if c.conn == nil {
		return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
	}

	c.conn = nil

	return nil
}

func getURL(o *Options) string {
	return fmt.Sprintf("%s:%s", o.Host, o.Port)
}
