package mrstorage

import (
	"context"
	"io"

	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// FileProvider - файловый провайдер.
	FileProvider interface {
		FileProviderConn
		FileProviderAPI
	}

	// FileProviderConn - управление открытым соединением файлового провайдера.
	FileProviderConn interface {
		Ping(ctx context.Context) error
		Close() error
	}

	// FileProviderAPI - файловый провайдер с возможностью загрузки, скачивания, удаления файла.
	FileProviderAPI interface {
		Info(ctx context.Context, filePath string) (mrtype.FileInfo, error)
		Download(ctx context.Context, filePath string) (mrtype.File, error)
		DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error)
		Upload(ctx context.Context, file mrtype.File) error
		Remove(ctx context.Context, filePath string) error
	}
)
