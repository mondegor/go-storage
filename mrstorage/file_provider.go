package mrstorage

import (
    "context"
    "io"
)

const (
    ModelNameFile = "File"
)

type (
    FileProvider interface {
        Download(ctx context.Context, file *File) error
        Upload(ctx context.Context, file *File) error
        Remove(ctx context.Context, filePath string) error
    }

    File struct {
        ContentType string
        Name string
        Size int64
        Body io.ReadCloser
    }
)
