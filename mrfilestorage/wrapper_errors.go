package mrfilestorage

import (
	"context"
	"io/fs"
	"os"

	"github.com/mondegor/go-sysmess/errors"
)

// wrapError - обёртывает ошибки файловой системы в стандартные ошибки приложения.
func (fp *FileProvider) wrapError(err error) error {
	if errors.Is(err, os.ErrNotExist) {
		return errors.ErrEventStorageNoRecordFound
	}

	if e := (*fs.PathError)(nil); errors.As(err, &e) {
		return errors.ErrInternalStorageQueryFailed.Wrap(err, "source_provider", providerName)
	}

	return errors.WrapInternalError(err, "failed", "source_provider", providerName)
}

// traceCmd - логирует выполняемую операцию для трассировки.
func (fp *FileProvider) traceCmd(ctx context.Context, command, filePath string) {
	fp.tracer.Trace(
		ctx,
		"source", providerName,
		"cmd", command,
		"file", fp.rootDir+filePath,
	)
}
