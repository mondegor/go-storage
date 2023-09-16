package mrfilestorage

import (
    "context"
    "io"
    "os"
    "path/filepath"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrlib"
)

func (c *nativeAdapter) Download(ctx context.Context, file *mrstorage.File) error {
    fullPath := c.getFullPath(file.Name)
    fi, err := os.Stat(fullPath)

    if err != nil {
        return err
    }

    fd, err := os.Open(fullPath)

    if err != nil {
        return err
    }

    file.ContentType = mrlib.MimeTypeByExt(filepath.Ext(file.Name))
    file.Size = fi.Size()
    file.Body = fd

    return nil
}

func (c *nativeAdapter) Upload(ctx context.Context, file *mrstorage.File) error {
    dst, err := os.Create(c.getFullPath(file.Name))

    if err != nil {
        return err
    }

    defer dst.Close()

    _, err = io.Copy(dst, file.Body)

    return err
}

func (c *nativeAdapter) Remove(ctx context.Context, filePath string) error {
    return os.Remove(c.getFullPath(filePath))
}

func (c *nativeAdapter) getFullPath(filePath string) string {
    return strings.Join([]string{c.rootDir, filePath}, "/")
}
