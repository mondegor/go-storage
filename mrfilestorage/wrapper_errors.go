package mrfilestorage

import (
	"context"
	"errors"
	"io/fs"
	"os"

	"github.com/mondegor/go-sysmess/mrerr/mr"
)

func (fp *FileProvider) wrapError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return mr.ErrStorageNoRowFound.Wrap(err)
	}

	var pathErr *fs.PathError
	if errors.As(err, &pathErr) {
		return mr.ErrStorageQueryFailed.Wrap(err)
	}

	return mr.ErrInternal.Wrap(err)
}

func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	fp.tracer.Trace(
		ctx,
		"source", providerName,
		"cmd", command,
		"file", fp.rootDir+filePath,
	)
}
