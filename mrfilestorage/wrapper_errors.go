package mrfilestorage

import (
	"errors"
	"io/fs"
	"os"

	"github.com/mondegor/go-webcore/mrcore"
)

func (fp *FileProvider) wrapError(err error, skip int) error {
	if errors.Is(err, os.ErrNotExist) {
		return mrcore.FactoryErrStorageNoRowFound.Caller(skip + 1).Wrap(err)
	}

	if _, ok := err.(*fs.PathError); ok {
		return mrcore.FactoryErrStorageQueryFailed.Caller(skip + 1).Wrap(err)
	}

	return mrcore.FactoryErrInternal.Caller(skip + 1).Wrap(err)
}
