package mrminio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func (fp *FileProvider) wrapError(err error) error {
	const skipFrame = 1

	minioErr, ok := err.(minio.ErrorResponse)
	if ok {
		// The specified key does not exist.
		if minioErr.Code == "NoSuchKey" {
			return mrcore.FactoryErrStorageNoRowFound.WithSkipFrame(skipFrame).Wrap(err)
		}

		return mrcore.FactoryErrStorageQueryFailed.WithSkipFrame(skipFrame).Wrap(err)
	}

	return mrcore.FactoryErrInternal.WithSkipFrame(skipFrame).Wrap(err)
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
