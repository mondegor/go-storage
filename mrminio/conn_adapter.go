package mrminio

import (
    "context"
    "fmt"
    "strings"

    "github.com/minio/minio-go/v7"
    "github.com/minio/minio-go/v7/pkg/credentials"
    "github.com/mondegor/go-webcore/mrcore"
)

// go get github.com/minio/minio-go/v7

const (
    connectionName = "minio"
)

type (
    connAdapter struct {
        conn *minio.Client
        backetName string
    }

    Options struct {
        Host string
        Port string
        UseSSL bool
        User string
        Password string
    }
)

func New(backetName string) *connAdapter {
    return &connAdapter{
        backetName: backetName,
    }
}

func (c *connAdapter) Connect(opt Options) error {
    if c.conn != nil {
        return mrcore.FactoryErrStorageConnectionIsAlreadyCreated.New(connectionName)
    }

    conn, err := minio.New(
        getUrl(&opt),
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

func (c *connAdapter) Ping(ctx context.Context) error {
    if c.conn == nil {
        return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
    }

    _, err := c.conn.GetBucketLocation(ctx, c.backetName)

    if err != nil && strings.Contains(err.Error(), "connection") {
        return mrcore.FactoryErrStorageConnectionFailed.Wrap(err, connectionName)
    }

    return nil
}

func (c *connAdapter) Cli() *minio.Client {
    return c.conn
}

func (c *connAdapter) Close() error {
    if c.conn == nil {
        return mrcore.FactoryErrStorageConnectionIsNotOpened.New(connectionName)
    }

    c.conn = nil

    return nil
}

func getUrl(o *Options) string {
    return fmt.Sprintf("%s:%s", o.Host, o.Port)
}
