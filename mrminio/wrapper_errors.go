package mrminio

import (
	"context"
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func (fp *FileProvider) wrapError(err error) error {
	var minioErr minio.ErrorResponse
	if errors.As(err, &minioErr) {
		// The specified key does not exist.
		if minioErr.Code == "NoSuchKey" {
			return mrcore.ErrStorageNoRowFound.Wrap(err)
		}

		return mrcore.ErrStorageQueryFailed.Wrap(err)
	}

	return mrcore.ErrInternal.Wrap(err)
}

func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", providerName).
		Str("cmd", command).
		Str("bucket", fp.bucketName).
		Str("file", filePath).
		Send()
}
