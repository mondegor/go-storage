package mrstorage

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAPI interface {
		Info(ctx context.Context, filePath string) (mrtype.FileInfo, error)
		Download(ctx context.Context, filePath string) (mrtype.File, error)
		Upload(ctx context.Context, file mrtype.File) error
		Remove(ctx context.Context, filePath string) error
	}
)
