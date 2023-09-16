package main

import (
    "context"

    "github.com/minio/minio-go/v7"
    "github.com/mondegor/go-storage/mrminio"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrtool"
)

func main() {
    logger := mrcore.DefaultLogger().With("mrminio")
    logger.Info("Create minio S3 connection")

    appHelper := mrtool.NewAppHelper(logger)

    opt := mrminio.Options{
        Host: "127.0.0.1",
        Port: "9000",
        UseSSL: false,
        User: "admin",
        Password: "12345678",
    }

    bucketName := "test-backet"

    minioAdapter := mrminio.New(bucketName)
    err := minioAdapter.Connect(opt)

    appHelper.ExitOnError(err)
    defer appHelper.Close(minioAdapter)

    appHelper.ExitOnError(minioAdapter.Ping(context.Background()))

    logger.Info("Create test backet")

    exists, err := minioAdapter.Cli().BucketExists(context.Background(), bucketName)
    appHelper.ExitOnError(err)

    if !exists {
        err = minioAdapter.Cli().MakeBucket(
            context.Background(),
            bucketName,
            minio.MakeBucketOptions{}, // "ru-central1"
        )
        appHelper.ExitOnError(err)

        logger.Info("Test backet created")
    }

    err = minioAdapter.Cli().RemoveBucket(context.Background(), bucketName)
    appHelper.ExitOnError(err)

    logger.Info("Test backet removed")
}
