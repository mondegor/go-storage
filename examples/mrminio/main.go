package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/litelog"
	"github.com/mondegor/go-sysmess/mrlog/slog"

	"github.com/mondegor/go-storage/mrminio"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

	l, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))
	logger := litelog.NewLogger(l)

	logger.Info("Create minio S3 connection")

	opts := mrminio.Options{
		Host:     "127.0.0.1",
		Port:     "9000",
		UseSSL:   false,
		User:     "admin",
		Password: "12345678",
	}

	ctx := context.Background()
	minioAdapter := mrminio.New(true, extfile.NewMimeTypeList([]extfile.MimeType{}), mrlog.NewDebugTracer(l))

	if err := minioAdapter.Connect(ctx, opts); err != nil {
		logger.Error("minioAdapter.Connect()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	defer minioAdapter.Close()

	if err := minioAdapter.Ping(ctx); err != nil {
		logger.Error("minioAdapter.Ping()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	logger.Info("Create test bucket")
	bucketName := "test-bucket"

	if created, err := minioAdapter.InitBucket(ctx, bucketName); err != nil {
		logger.Error("minioAdapter.InitBucket()", "error", err) // mrlog.Fatal
		os.Exit(1)
	} else {
		if created {
			logger.Info("Bucket created", "bucket", bucketName)
		} else {
			logger.Info("Bucket exists, OK", "bucket", bucketName)
		}
	}

	minioCli, err := minioAdapter.Cli()
	if err != nil {
		logger.Error("minioAdapter.Cli()", "error", err) // mrlog.Fatal
		os.Exit(1)
	}

	if err := minioCli.RemoveBucket(ctx, bucketName); err != nil {
		logger.Error("minioCli.RemoveBucket()", "error", err) // mrlog.Fatal
		os.Exit(1)
	} else {
		logger.Info("Test bucket removed")
	}
}
