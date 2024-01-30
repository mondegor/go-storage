package main

import (
	"context"

	"github.com/mondegor/go-storage/mrminio"
	"github.com/mondegor/go-webcore/mrlog"
)

func main() {
	logger := mrlog.New(mrlog.TraceLevel)
	ctx := mrlog.WithContext(context.Background(), logger)

	logger.Info().Msg("Create minio S3 connection")

	opts := mrminio.Options{
		Host:     "127.0.0.1",
		Port:     "9000",
		UseSSL:   false,
		User:     "admin",
		Password: "12345678",
	}

	minioAdapter := mrminio.New(true)

	if err := minioAdapter.Connect(ctx, opts); err != nil {
		logger.Fatal().Err(err).Msg("minioAdapter.Connect() error")
	}

	defer minioAdapter.Close()

	if err := minioAdapter.Ping(ctx); err != nil {
		logger.Fatal().Err(err).Msg("minioAdapter.Ping() error")
	}

	logger.Info().Msg("Create test bucket")
	bucketName := "test-bucket"

	if created, err := minioAdapter.InitBucket(ctx, bucketName); err != nil {
		logger.Fatal().Err(err).Msg("minioAdapter.InitBucket() error")
	} else {
		if created {
			logger.Info().Msgf("Bucket '%s' created", bucketName)
		} else {
			logger.Info().Msgf("Bucket '%s' exists, OK", bucketName)
		}
	}

	if err := minioAdapter.Cli().RemoveBucket(ctx, bucketName); err != nil {
		logger.Fatal().Err(err).Msg("minioAdapter.Cli().RemoveBucket() error")
	} else {
		logger.Info().Msg("Test bucket removed")
	}
}
