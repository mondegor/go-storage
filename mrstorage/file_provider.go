package mrstorage

import (
	"context"
	"io"

	modelmedia "github.com/mondegor/go-sysmess/mrmodel/media"
)

type (
	// FileProvider - файловый провайдер для работы с хранилищами файлов.
	// Объединяет управление соединением и API для операций с файлами.
	FileProvider interface {
		FileProviderConn
		FileProviderAPI
	}

	// FileProviderConn - управление подключением файлового провайдера.
	FileProviderConn interface {
		Ping(ctx context.Context) error
		Close() error
	}

	// FileProviderAPI - API файлового провайдера для операций с файлами.
	// Позволяет получать информацию, загружать, сохранять и удалять файлы.
	FileProviderAPI interface {
		Info(ctx context.Context, filePath string) (modelmedia.FileInfo, error)
		Download(ctx context.Context, filePath string) (modelmedia.File, error)
		DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error)
		Upload(ctx context.Context, file modelmedia.File) error
		Remove(ctx context.Context, filePath string) error
	}
)
