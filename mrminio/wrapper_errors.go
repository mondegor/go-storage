package mrminio

import (
	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrcore"
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
