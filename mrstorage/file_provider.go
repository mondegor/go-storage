package mrstorage

import (
    "context"

    "github.com/mondegor/go-webcore/mrtype"
)

type (
    FileProviderAPI interface {
        Info(ctx context.Context, path string) (mrtype.FileInfo, error)
        Download(ctx context.Context, path string) (*mrtype.File, error)
        // Downloads(ctx context.Context, path string) (*mrtype.DownloadedFile, error) // ListObjects :TODO: получение списка объектов
        Upload(ctx context.Context, file *mrtype.File) error
        Remove(ctx context.Context, path string) error
    }
)
