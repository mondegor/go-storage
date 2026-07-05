package mrminio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-sysmess/errors"
)

// wrapError - обёртывает ошибки MinIO в стандартные ошибки приложения.
func (fp *FileProvider) wrapError(err error) error {
	if e := (*minio.ErrorResponse)(nil); errors.As(err, &e) {
		// The specified key does not exist.
		if e.Code == "NoSuchKey" {
			return errors.ErrEventStorageNoRecordFound
		}

		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source_provider", providerName)
	}

	return errors.WrapInternalError(err, "failed", "source_provider", providerName)
}

// traceCmd - логирует выполняемую операцию для трассировки.
func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	fp.tracer.Trace(
		ctx,
		"source", providerName,
		"cmd", command,
		"bucket", fp.bucketName,
		"file", filePath,
	)
}
