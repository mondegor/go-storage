package mrminio

import (
	"context"
	"errors"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-sysmess/mrerr/mr"
)

func (fp *FileProvider) wrapError(err error) error {
	var minioErr minio.ErrorResponse
	if errors.As(err, &minioErr) {
		// The specified key does not exist.
		if minioErr.Code == "NoSuchKey" {
			return mr.ErrStorageNoRowFound.Wrap(err)
		}

		return mr.ErrStorageQueryFailed.Wrap(err)
	}

	return mr.ErrInternal.Wrap(err)
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
