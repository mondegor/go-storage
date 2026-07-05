package mrminio

import (
	"context"
	"fmt"
	"io"
	"path"

	"github.com/minio/minio-go/v7"
	"github.com/mondegor/go-core/errors"
	modelmedia "github.com/mondegor/go-core/mrmodel/media"
	"github.com/mondegor/go-core/util/casttype"
)

// https://min.io/docs/minio/linux/developers/go/API.html

const (
	// providerName - имя провайдера для логирования и трассировки.
	providerName = "Minio"
)

type (
	// FileProvider - файловый провайдер для S3-совместимого хранилища MinIO.
	// Позволяет читать, сохранять, удалять файлы и проверять их метаинформацию.
	FileProvider struct {
		*ConnAdapter
		bucketName string
	}
)

// NewFileProvider - создаёт объект FileProvider для работы с конкретным бакетом.
func NewFileProvider(conn *ConnAdapter, bucketName string) *FileProvider {
	return &FileProvider{
		ConnAdapter: conn,
		bucketName:  bucketName,
	}
}

// Info - возвращает метаинформацию о файле (размер, тип контента, даты).
func (fp *FileProvider) Info(ctx context.Context, filePath string) (modelmedia.FileInfo, error) {
	fp.traceCmd(ctx, "Info", filePath)

	info, err := fp.conn.StatObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.StatObjectOptions{},
	)
	if err != nil {
		return modelmedia.FileInfo{}, fp.wrapError(err)
	}

	fileInfo, err := fp.getFileInfo(&info)
	if err != nil {
		return modelmedia.FileInfo{}, fp.wrapError(err)
	}

	return fileInfo, nil
}

// Download - открывает файл и возвращает его метаинформацию вместе с содержимым.
// Возвращаемая структура mrmodel.File включает тело файла (Body), которое нужно закрыть после использования.
func (fp *FileProvider) Download(ctx context.Context, filePath string) (modelmedia.File, error) {
	fp.traceCmd(ctx, "Download", filePath)

	object, err := fp.openObject(ctx, filePath)
	if err != nil {
		return modelmedia.File{}, fp.wrapError(err)
	}

	info, err := object.Stat()
	if err != nil {
		_ = object.Close()

		return modelmedia.File{}, fp.wrapError(err)
	}

	fileInfo, err := fp.getFileInfo(&info)
	if err != nil {
		return modelmedia.File{}, fp.wrapError(err)
	}

	return modelmedia.File{
		FileInfo: fileInfo,
		Body:     object,
	}, nil
}

// DownloadFile - открывает файл и возвращает только его содержимое как io.ReadCloser.
// В отличие от Download, не включает метаинформацию.
// Вызывающая сторона обязана закрыть ReadCloser после использования.
func (fp *FileProvider) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fp.traceCmd(ctx, "DownloadFile", filePath)

	object, err := fp.openObject(ctx, filePath)
	if err != nil { // :TODO: ошибки нет даже если filePath не найден
		return nil, fp.wrapError(err)
	}

	if _, err = object.Stat(); err != nil {
		_ = object.Close()

		return nil, fp.wrapError(err)
	}

	return object, nil
}

// Upload - сохраняет файл в бакет MinIO.
// Если размер файла не указан (равен 0), автоматически вычисляет его при загрузке.
// Определяет ContentType по расширению файла, если он не указан явно.
// Устанавливает Content-Disposition для корректного скачивания с оригинальным именем.
func (fp *FileProvider) Upload(ctx context.Context, file modelmedia.File) error {
	fp.traceCmd(ctx, "Upload", file.Path)

	fileSize := file.Size

	if fileSize == 0 {
		fileSize = -1 // -1 - calculating size
	}

	contentType, err := fp.getContentType(file.ContentType, file.Path)
	if err != nil {
		return fp.wrapError(err)
	}

	_, err = fp.conn.PutObject(
		ctx,
		fp.bucketName,
		file.Path,
		file.Body,
		fileSize,
		minio.PutObjectOptions{
			ContentType:        contentType,
			ContentDisposition: fp.getContentDisposition(file.OriginalName),
		},
	)
	if err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Remove - удаляет файл из бакета MinIO.
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

// openObject - открывает объект из бакета MinIO для чтения.
func (fp *FileProvider) openObject(ctx context.Context, filePath string) (*minio.Object, error) {
	return fp.conn.GetObject(
		ctx,
		fp.bucketName,
		filePath,
		minio.GetObjectOptions{},
	)
}

// getFileInfo - извлекает метаинформацию из minio.ObjectInfo.
// Определяет ContentType, парсит оригинальное имя файла из Content-Disposition.
func (fp *FileProvider) getFileInfo(info *minio.ObjectInfo) (modelmedia.FileInfo, error) {
	contentType, err := fp.getContentType(info.ContentType, info.Key)
	if err != nil {
		return modelmedia.FileInfo{}, err
	}

	if info.Size < 0 {
		return modelmedia.FileInfo{}, errors.NewInternalError(
			"file size is negative",
			"file", info.Key,
		)
	}

	return modelmedia.FileInfo{
		ContentType:  contentType,
		OriginalName: fp.parseOriginalName(info.Metadata.Get("Content-Disposition")),
		Name:         path.Base(info.Key),
		Path:         info.Key,
		Size:         info.Size,
		UpdatedAt:    casttype.TimeToPointer(info.LastModified),
	}, nil
}

// getContentDisposition - формирует значение заголовка Content-Disposition для скачивания файла.
// Если оригинальное имя файла пустое, возвращает пустую строку.
// :TODO: добавить экранирование значения для безопасности.
func (fp *FileProvider) getContentDisposition(value string) string {
	if value == "" {
		return ""
	}

	return fmt.Sprintf("attachment; filename=\"%s\"", value) // :TODO: escape value
}

// getContentType - определяет тип контента файла.
// Если ContentType уже указан, возвращает его.
// Иначе определяет по расширению файла через список MIME-типов.
func (fp *FileProvider) getContentType(contentType, fileName string) (string, error) {
	if contentType != "" {
		return contentType, nil
	}

	contentType, err := fp.mimeTypes.ContentTypeByExt(path.Ext(fileName))
	if err != nil {
		return "", err
	}

	return contentType, nil
}

// parseOriginalName - извлекает оригинальное имя файла из заголовка Content-Disposition.
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
