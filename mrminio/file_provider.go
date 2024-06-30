package mrminio

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-webcore/mrtype"
)

// https://min.io/docs/minio/linux/developers/go/API.html

const (
	providerName = "Minio"
)

type (
	// FileProvider - файловый провайдер S3,
	// позволяет читать, сохранять, удалять файлы.
	FileProvider struct {
		*ConnAdapter
		bucketName string
	}
)

// NewFileProvider - создаёт объект FileProvider.
func NewFileProvider(conn *ConnAdapter, bucketName string) *FileProvider {
	return &FileProvider{
		ConnAdapter: conn,
		bucketName:  bucketName,
	}
}

// Info - comment method.
func (fp *FileProvider) Info(ctx context.Context, filePath string) (mrtype.FileInfo, error) {
	fp.traceCmd(ctx, "Info", filePath)

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

// Download - comment method.
func (fp *FileProvider) Download(ctx context.Context, filePath string) (mrtype.File, error) {
	fp.traceCmd(ctx, "Download", filePath)

	object, err := fp.openObject(ctx, filePath)
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

// DownloadFile - comment method.
func (fp *FileProvider) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fp.traceCmd(ctx, "DownloadFile", filePath)

	object, err := fp.openObject(ctx, filePath)
	if err != nil { // :TODO: ошибки нет даже если filePath не найден
		return nil, fp.wrapError(err)
	}

	if _, err = object.Stat(); err != nil {
		object.Close()

		return nil, fp.wrapError(err)
	}

	return object, nil
}

// Upload - comment method.
func (fp *FileProvider) Upload(ctx context.Context, file mrtype.File) error {
	fp.traceCmd(ctx, "Upload", file.Path)

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
	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Remove - comment method.
func (fp *FileProvider) Remove(ctx context.Context, filePath string) error {
	fp.traceCmd(ctx, "Remove", filePath)

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

func (fp *FileProvider) openObject(ctx context.Context, filePath string) (*minio.Object, error) {
	return fp.conn.GetObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.GetObjectOptions{},
	)
}

func (fp *FileProvider) getFileInfo(info *minio.ObjectInfo) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType:  fp.getContentType(info.ContentType, info.Key),
		OriginalName: fp.parseOriginalName(info.Metadata.Get("Content-Disposition")),
		Name:         path.Base(info.Key),
		Path:         info.Key,
		Size:         info.Size,
		UpdatedAt:    mrtype.TimeToPointer(info.LastModified),
	}
}

func (fp *FileProvider) getContentDisposition(value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("attachment; filename=\"%s\"", value) // :TODO: escape value
}

func (fp *FileProvider) getContentType(contentType, fileName string) string {
	if contentType != "" {
		return contentType
	}

	return fp.mimeTypes.ContentTypeByFileName(fileName)
}

func (fp *FileProvider) parseOriginalName(contentDisposition string) string {
	const (
		prefix    = "attachment; filename=\""
		minLength = 23 // len of prefix + '"'
	)

	length := len(contentDisposition)

	if length > minLength &&
		contentDisposition[:minLength-1] == prefix &&
		contentDisposition[length-1] == '"' {
		return contentDisposition[minLength-1 : length-1]
	}

	return contentDisposition
}
