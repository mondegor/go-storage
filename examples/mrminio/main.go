package main

import (
	"context"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtool"
)

func main() {
	logger := mrcore.Log().With("mrminio")
	logger.Info("Create minio S3 connection")

	appHelper := mrtool.NewAppHelper(logger)

	opt := mrminio.Options{
		Host:     "127.0.0.1",
		Port:     "9000",
		UseSSL:   false,
		User:     "admin",
		Password: "12345678",
	}

	minioAdapter := mrminio.New()
	err := minioAdapter.Connect(opt)

	appHelper.ExitOnError(err)
	defer appHelper.Close(minioAdapter)

	appHelper.ExitOnError(minioAdapter.Ping(context.Background()))

	logger.Info("Create test bucket")

	bucketName := "test-bucket"

	created, err := minioAdapter.InitBucket(context.Background(), bucketName, true)
	appHelper.ExitOnError(err)

	if created {
		mrcore.LogInfo("Bucket '%s' created", bucketName)
	} else {
		mrcore.LogInfo("Bucket '%s' exists, OK", bucketName)
	}

	err = minioAdapter.Cli().RemoveBucket(context.Background(), bucketName)
	appHelper.ExitOnError(err)

	logger.Info("Test bucket removed")
}
