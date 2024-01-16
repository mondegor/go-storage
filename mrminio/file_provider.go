package mrminio

import (
	"context"
	"fmt"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

// https://min.io/docs/minio/linux/developers/go/API.html

const (
	providerName = "Minio"
)

type (
	FileProvider struct {
		*ConnAdapter
		bucketName string
	}
)

func NewFileProvider(conn *ConnAdapter, bucketName string) *FileProvider {
	return &FileProvider{
		ConnAdapter: conn,
		bucketName:  bucketName,
	}
}

func (fp *FileProvider) Info(ctx context.Context, filePath string) (mrtype.FileInfo, error) {
	fp.debugCmd(ctx, "Info", filePath)

	info, err := fp.conn.StatObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.StatObjectOptions{},
	)

	if err != nil {
		return mrtype.FileInfo{}, fp.wrapError(err)
	}

	return fp.getFileInfo(&info), nil
}

func (fp *FileProvider) Download(ctx context.Context, filePath string) (mrtype.File, error) {
	fp.debugCmd(ctx, "Download", filePath)

	object, err := fp.conn.GetObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.GetObjectOptions{},
	)

	if err != nil {
		return mrtype.File{}, fp.wrapError(err)
	}

	info, err := object.Stat()

	if err != nil {
		object.Close()
		return mrtype.File{}, fp.wrapError(err)
	}

	return mrtype.File{
		FileInfo: fp.getFileInfo(&info),
		Body:     object,
	}, nil
}

func (fp *FileProvider) Upload(ctx context.Context, file mrtype.File) error {
	fp.debugCmd(ctx, "Upload", file.Path)

	_, err := fp.conn.PutObject(
		ctx,
		fp.bucketName,
		file.Path,
		file.Body,
		file.Size, // -1 - calculating size
		minio.PutObjectOptions{
			ContentType:        mrlib.MimeType(file.ContentType, file.Path),
			ContentDisposition: fp.getContentDisposition(file.OriginalName),
		},
	)

	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

func (fp *FileProvider) Remove(ctx context.Context, filePath string) error {
	fp.debugCmd(ctx, "Remove", filePath)

	err := fp.conn.RemoveObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.RemoveObjectOptions{},
	)

	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

func (fp *FileProvider) getFileInfo(info *minio.ObjectInfo) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType:  mrlib.MimeType(info.ContentType, info.Key),
		OriginalName: fp.parseOriginalName(info.Metadata.Get("Content-Disposition")),
		Name:         path.Base(info.Key),
		Path:         info.Key,
		Size:         info.Size,
		ModifiedAt:   mrtype.TimePointer(info.LastModified),
	}
}

func (fp *FileProvider) getContentDisposition(value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("attachment; filename=\"%s\"", value) // :TODO: escape value
}

func (fp *FileProvider) parseOriginalName(contentDisposition string) string {
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
