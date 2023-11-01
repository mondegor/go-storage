package mrstorage

import (
    "context"

    "github.com/mondegor/go-storage/mrentity"
)

type (
    FileProvider interface {
        Download(ctx context.Context, file *mrentity.File) error
        Upload(ctx context.Context, file *mrentity.File) error
        Remove(ctx context.Context, filePath string) error
    }
)
