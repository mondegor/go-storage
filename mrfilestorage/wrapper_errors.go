package mrfilestorage

import (
	"context"
	"io/fs"
	"os"

	"github.com/mondegor/go-sysmess/errors"
)

func (fp *FileProvider) wrapError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return errors.ErrEventStorageNoRowFound
	}

	if e := (*fs.PathError)(nil); errors.As(err, &e) {
		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source_provider", providerName)
	}

	return errors.WrapInternalError(err, "failed", "source_provider", providerName)
}

func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	fp.tracer.Trace(
		ctx,
		"source", providerName,
		"cmd", command,
		"file", fp.rootDir+filePath,
	)
}
