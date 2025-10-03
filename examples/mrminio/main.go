package main

import (
	"context"
	"os"

	"github.com/mondegor/go-sysmess/mrerr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrlog/slog"
	"github.com/mondegor/go-sysmess/mrtrace/logtracer"

	"github.com/mondegor/go-storage/mrminio"
)

func main() {
	mrerr.InitDefaultOptions(mrerr.DefaultOptionsHandler())

	logger, _ := slog.NewLoggerAdapter(slog.WithWriter(os.Stdout))

	mrlog.Info(logger, "Create minio S3 connection")

	opts := mrminio.Options{
		Host:     "127.0.0.1",
		Port:     "9000",
		UseSSL:   false,
		User:     "admin",
		Password: "12345678",
	}

	ctx := context.Background()
	minioAdapter := mrminio.New(true, extfile.NewMimeTypeList([]extfile.MimeType{}), logtracer.New(logger))

	if err := minioAdapter.Connect(ctx, opts); err != nil {
		mrlog.Fatal(logger, "minioAdapter.Connect()", "error", err) // mrlog.Fatal
	}

	defer minioAdapter.Close()

	if err := minioAdapter.Ping(ctx); err != nil {
		mrlog.Fatal(logger, "minioAdapter.Ping()", "error", err) // mrlog.Fatal
	}

	mrlog.Info(logger, "Create test bucket")
	bucketName := "test-bucket"

	if created, err := minioAdapter.InitBucket(ctx, bucketName); err != nil {
		mrlog.Fatal(logger, "minioAdapter.InitBucket()", "error", err) // mrlog.Fatal
	} else {
		if created {
			mrlog.Info(logger, "Bucket created", "bucket", bucketName)
		} else {
			mrlog.Info(logger, "Bucket exists, OK", "bucket", bucketName)
		}
	}

	minioCli, err := minioAdapter.Cli()
	if err != nil {
		mrlog.Fatal(logger, "minioAdapter.Cli()", "error", err) // mrlog.Fatal
	}

	if err := minioCli.RemoveBucket(ctx, bucketName); err != nil {
		mrlog.Fatal(logger, "minioCli.RemoveBucket()", "error", err) // mrlog.Fatal
	} else {
		mrlog.Info(logger, "Test bucket removed")
	}
}
