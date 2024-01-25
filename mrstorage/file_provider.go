package mrstorage

import (
	"context"
	"io"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAPI interface {
		Info(ctx context.Context, filePath string) (mrtype.FileInfo, error)
		Download(ctx context.Context, filePath string) (mrtype.File, error)
		DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error)
		Upload(ctx context.Context, file mrtype.File) error
		Remove(ctx context.Context, filePath string) error
	}
)
