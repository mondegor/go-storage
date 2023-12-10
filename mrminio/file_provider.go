package mrminio

import (
	"context"
	"fmt"
	"path/filepath"
	"strings"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

// https://min.io/docs/minio/linux/developers/go/API.html

type (
	FileProvider struct {
		*ConnAdapter
		bucketName string
		baseDir    string
	}
)

func NewFileProvider(conn *ConnAdapter, bucketName string) *FileProvider {
	return &FileProvider{
		ConnAdapter: conn,
		bucketName:  bucketName,
	}
}

func (fp *FileProvider) WithBaseDir(value string) (mrstorage.ExtFileProviderAPI, error) {
	if value != "" {
		value = strings.Trim(value, "/")

		if value != "" {
			value += "/"
		}
	}

	if fp.baseDir == value {
		return fp, nil
	}

	c := *fp
	c.baseDir = value

	return &c, nil
}

func (fp *FileProvider) Info(ctx context.Context, fileName string) (mrtype.FileInfo, error) {
	info, err := fp.conn.StatObject(
		ctx,
		fp.bucketName,
		fp.baseDir+fileName,
		minio.StatObjectOptions{},
	)

	if err != nil {
		return mrtype.FileInfo{}, fp.wrapError(err)
	}

	return fp.getFileInfo(&info, fileName), nil
}

func (fp *FileProvider) Download(ctx context.Context, fileName string) (*mrtype.File, error) {
	object, err := fp.conn.GetObject(
		ctx,
		fp.bucketName,
		fp.baseDir+fileName,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return nil, fp.wrapError(err)
	}

	info, err := object.Stat()

	if err != nil {
		object.Close()
		return nil, fp.wrapError(err)
	}

	return &mrtype.File{
		FileInfo: fp.getFileInfo(&info, fileName),
		Body:     object,
	}, nil
}

func (fp *FileProvider) Upload(ctx context.Context, file *mrtype.File) error {
	_, err := fp.conn.PutObject(
		ctx,
		fp.bucketName,
		fp.baseDir+file.Path,
		file.Body,
		file.Size, // -1 - calculating size
		minio.PutObjectOptions{
			ContentType:        fp.getContentType(file.ContentType, file.Path),
			ContentDisposition: fp.getContentDisposition(file.OriginalName),
		},
	)

	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

func (fp *FileProvider) Remove(ctx context.Context, fileName string) error {
	err := fp.conn.RemoveObject(
		ctx,
		fp.bucketName,
		fp.baseDir+fileName,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

func (fp *FileProvider) getFileInfo(info *minio.ObjectInfo, fileName string) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType:  fp.getContentType(info.ContentType, fileName),
		OriginalName: fp.getOriginalName(info.Metadata.Get("Content-Disposition")),
		Name:         info.Key,
		LastModified: info.LastModified,
		Size:         info.Size,
	}
}

func (fp *FileProvider) getContentType(value, fileName string) string {
	if value != "" {
		return value
	}

	return mrlib.MimeTypeByExt(filepath.Ext(fileName))
}

func (fp *FileProvider) getContentDisposition(value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("attachment; filename=\"%s\"", value) // :TODO: escape value
}

func (fp *FileProvider) getOriginalName(contentDisposition string) string {
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
