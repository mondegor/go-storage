package mrstorage

import (
	"context"

	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProviderAPI interface {
		Info(ctx context.Context, fileName string) (mrtype.FileInfo, error)
		Download(ctx context.Context, fileName string) (*mrtype.File, error)
		// Downloads(ctx context.Context, fileName string) (*mrtype.DownloadedFile, error) // ListObjects :TODO: получение списка объектов
		Upload(ctx context.Context, file *mrtype.File) error
		Remove(ctx context.Context, fileName string) error
	}

	// ExtFileProviderAPI - WARNING: use only when starting the main process
	ExtFileProviderAPI interface {
		WithBaseDir(value string) (ExtFileProviderAPI, error)
		FileProviderAPI
	}
)
