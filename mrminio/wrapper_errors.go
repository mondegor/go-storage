package mrminio

import (
	"context"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

func (fp *FileProvider) wrapError(err error) error {
	minioErr, ok := err.(minio.ErrorResponse)

	if ok {
		// The specified key does not exist.
		if minioErr.Code == "NoSuchKey" {
			return mrcore.FactoryErrStorageNoRowFound.Caller(1).Wrap(err)
		}

		return mrcore.FactoryErrStorageQueryFailed.Caller(1).Wrap(err)
	}

	return mrcore.FactoryErrInternal.Caller(1).Wrap(err)
}

func (fp *FileProvider) debugCmd(ctx context.Context, command, filePath string) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, bucket=%s, file=%s",
		providerName,
		command,
		fp.bucketName,
		filePath,
	)
}
