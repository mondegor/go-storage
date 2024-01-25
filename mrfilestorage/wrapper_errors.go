package mrfilestorage

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrctx"
)

func (fp *FileProvider) wrapError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return mrcore.FactoryErrStorageNoRowFound.Caller(1).Wrap(err)
	}

	if _, ok := err.(*fs.PathError); ok {
		return mrcore.FactoryErrStorageQueryFailed.Caller(1).Wrap(err)
	}

	return mrcore.FactoryErrInternal.Caller(1).Wrap(err)
}

func (fp *FileProvider) debugCmd(ctx context.Context, command, filePath string) {
	mrctx.Logger(ctx).Debug(
		"%s: cmd=%s, file=%s",
		providerName,
		command,
		fp.rootDir+filePath,
	)
}
