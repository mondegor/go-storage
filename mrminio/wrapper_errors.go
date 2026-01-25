package mrminio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-sysmess/errors"
)

func (fp *FileProvider) wrapError(err error) error {
	if e := (*minio.ErrorResponse)(nil); errors.As(err, &e) {
		// The specified key does not exist.
		if e.Code == "NoSuchKey" {
			return errors.ErrEventStorageNoRowFound
		}

		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source_provider", providerName)
	}

	return errors.WrapInternalError(err, "failed", "source_provider", providerName)
}

func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	fp.tracer.Trace(
		ctx,
		"source", providerName,
		"cmd", command,
		"bucket", fp.bucketName,
		"file", filePath,
	)
}
