package mrminio

import (
	"context"
	"fmt"
	"path/filepath"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

// https://min.io/docs/minio/linux/developers/go/API.html

type (
	fileProvider struct {
		*ConnAdapter
		bucketName string
	}
)

func NewFileProvider(conn *ConnAdapter, bucketName string) *fileProvider {
	return &fileProvider{
		ConnAdapter: conn,
		bucketName:  bucketName,
	}
}

func (fp *fileProvider) Info(ctx context.Context, path string) (mrtype.FileInfo, error) {
	info, err := fp.conn.StatObject(
		ctx,
		fp.bucketName,
		path,
		minio.StatObjectOptions{},
	)

	if err != nil {
		return mrtype.FileInfo{}, err
	}

	return fp.getFileInfo(&info, path), nil
}

func (fp *fileProvider) Download(ctx context.Context, path string) (*mrtype.File, error) {
	object, err := fp.conn.GetObject(
		ctx,
		fp.bucketName,
		path,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, err
	}

	info, err := object.Stat()

	if err != nil {
		object.Close()
		return nil, err
	}

	return &mrtype.File{
		FileInfo: fp.getFileInfo(&info, path),
		Body:     object,
	}, nil
}

func (fp *fileProvider) Upload(ctx context.Context, file *mrtype.File) error {
	_, err := fp.conn.PutObject(
		ctx,
		fp.bucketName,
		file.Path,
		file.Body,
		file.Size, // -1 - calculating size
		minio.PutObjectOptions{
			ContentType:        fp.getContentType(file.ContentType, file.Path),
			ContentDisposition: fp.getContentDisposition(file.OriginalName),
		},
	)

	return err
}

func (fp *fileProvider) Remove(ctx context.Context, path string) error {
	return fp.conn.RemoveObject(
		ctx,
		fp.bucketName,
		path,
		minio.RemoveObjectOptions{},
	)
}

func (fp *fileProvider) getFileInfo(info *minio.ObjectInfo, path string) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType:  fp.getContentType(info.ContentType, path),
		OriginalName: fp.getOriginalName(info.Metadata.Get("Content-Disposition")),
		Name:         info.Key,
		LastModified: info.LastModified,
		Size:         info.Size,
	}
}

func (fp *fileProvider) getContentType(value, path string) string {
	if value != "" {
		return value
	}

	return mrlib.MimeTypeByExt(filepath.Ext(path))
}

func (fp *fileProvider) getContentDisposition(value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("attachment; filename=\"%s\"", value) // :TODO: escape value
}

func (fp *fileProvider) getOriginalName(contentDisposition string) string {
	const prefix = "attachment; filename=\""
	const minLength = 23 // len of prefix + '"'

	length := len(contentDisposition)

	if length > minLength &&
		contentDisposition[:minLength-1] == prefix &&
		contentDisposition[length-1] == '"' {
		return contentDisposition[minLength-1 : length-1]
	}

	return contentDisposition
}
