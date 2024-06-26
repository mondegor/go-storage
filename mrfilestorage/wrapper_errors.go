package mrfilestorage

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
)

func (fp *FileProvider) wrapError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return mrcore.ErrStorageNoRowFound.Wrap(err)
	}

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return mrcore.ErrStorageQueryFailed.Wrap(err)
	}

	return mrcore.ErrInternal.Wrap(err)
}

func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	mrlog.Ctx(ctx).
		Trace().
		Str("source", providerName).
		Str("cmd", command).
		Str("file", fp.rootDir+filePath).
		Send()
}
