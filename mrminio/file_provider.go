package mrminio

import (
    "context"
    "path/filepath"

    "github.com/minio/minio-go/v7"
    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-webcore/mrlib"
)

func (c *connAdapter) Download(ctx context.Context, file *mrentity.File) error {
    reader, err := c.conn.GetObject(
        ctx,
        c.backetName,
        file.Name,
        minio.GetObjectOptions{},
    )

    if err != nil {
        return err
    }

    fi, err := reader.Stat()

    if err != nil {
        reader.Close()
        return err
    }

    if fi.ContentType != "" {
        file.ContentType = fi.ContentType
    } else {
        file.ContentType = mrlib.MimeTypeByExt(filepath.Ext(file.Name))
    }

    file.Size = fi.Size
    file.Body = reader

    return nil
}

func (c *connAdapter) Upload(ctx context.Context, file *mrentity.File) error {
    var opts minio.PutObjectOptions

    if file.ContentType != "" {
        opts.ContentType = file.ContentType
    } else {
        opts.ContentType = mrlib.MimeTypeByExt(filepath.Ext(file.Name))
    }

    _, err := c.conn.PutObject(
        ctx,
        c.backetName,
        file.Name,
        file.Body,
        file.Size,
        opts,
    )

    return err
}

func (c *connAdapter) Remove(ctx context.Context, filePath string) error {
    return c.conn.RemoveObject(
        ctx,
        c.backetName,
        filePath,
        minio.RemoveObjectOptions{},
    )
}
