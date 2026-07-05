package mrfilestorage

import (
	"context"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/mondegor/go-sysmess/errors"
	modelmedia "github.com/mondegor/go-sysmess/mrmodel/media"
	"github.com/mondegor/go-sysmess/mrtrace"
	"github.com/mondegor/go-sysmess/util/casttype"
)

const (
	// providerName - имя провайдера для логирования и трассировки.
	providerName = "FileStorage"

	// testFile - имя временного файла для проверки работоспособности хранилища в методе Ping.
	testFile = "testFile-d6b6943c-e1b2-4625-b133-9805a5cf5f8d"
)

type (
	// FileProvider - файловый провайдер, работающий с нативной файловой системой.
	// Позволяет читать, сохранять, удалять файлы и проверять работоспособность хранилища.
	FileProvider struct {
		fs      *FileSystem    // fs - объект для работы с файловой системой
		tracer  mrtrace.Tracer // tracer - трассировщик для логирования операций
		rootDir string         // rootDir - корневая директория хранилища (всегда заканчивается на "/")
		muPing  sync.Mutex     // muPing - мьютекс для синхронизации метода Ping
	}
)

// NewFileProvider - создаёт объект FileProvider.
func NewFileProvider(fs *FileSystem, tracer mrtrace.Tracer, rootDir string) *FileProvider {
	return &FileProvider{
		fs:      fs,
		tracer:  tracer,
		rootDir: strings.TrimRight(rootDir, "/") + "/",
	}
}

// Info - возвращает метаинформацию о файле (размер, тип контента, даты).
func (fp *FileProvider) Info(ctx context.Context, filePath string) (modelmedia.FileInfo, error) {
	fp.traceCmd(ctx, "Info", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return modelmedia.FileInfo{}, err
	}

	fi, err := os.Stat(fp.rootDir + filePath)
	if err != nil {
		return modelmedia.FileInfo{}, fp.wrapError(err)
	}

	fileInfo, err := fp.getFileInfo(filePath, fi)
	if err != nil {
		return modelmedia.FileInfo{}, fp.wrapError(err)
	}

	return fileInfo, nil
}

// Download - открывает файл и возвращает его метаинформацию вместе с содержимым.
// Возвращаемая структура modelmedia.File включает тело файла (Body), которое нужно закрыть после использования.
func (fp *FileProvider) Download(ctx context.Context, filePath string) (modelmedia.File, error) {
	fp.traceCmd(ctx, "Download", filePath)

	fd, err := fp.openFile(ctx, filePath)
	if err != nil {
		return modelmedia.File{}, fp.wrapError(err)
	}

	fi, err := fd.Stat()
	if err != nil {
		return modelmedia.File{}, fp.wrapError(err)
	}

	fileInfo, err := fp.getFileInfo(filePath, fi)
	if err != nil {
		return modelmedia.File{}, fp.wrapError(err)
	}

	return modelmedia.File{
		FileInfo: fileInfo,
		Body:     fd,
	}, nil
}

// DownloadFile - открывает файл и возвращает только его содержимое как io.ReadCloser.
// В отличие от Download, не включает метаинформацию.
// Вызывающая сторона обязана закрыть ReadCloser после использования.
func (fp *FileProvider) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fp.traceCmd(ctx, "DownloadContent", filePath)

	fd, err := fp.openFile(ctx, filePath)
	if err != nil {
		return nil, fp.wrapError(err)
	}

	return fd, nil
}

// Upload - сохраняет файл в хранилище.
// Автоматически создаёт необходимые директории, если они не существуют.
func (fp *FileProvider) Upload(ctx context.Context, file modelmedia.File) error {
	fp.traceCmd(ctx, "Upload", file.Path)

	if err := fp.checkFilePath(file.Path); err != nil {
		return err
	}

	if dirPath := path.Dir(file.Path); dirPath != "" {
		if err := fp.fs.CreateDirIfNotExists(fp.rootDir, dirPath); err != nil {
			return fp.wrapError(err)
		}
	}

	dst, err := os.Create(fp.rootDir + file.Path)
	if err != nil {
		return fp.wrapError(err)
	}

	defer func() {
		_ = dst.Close()
	}()

	if _, err = io.Copy(dst, file.Body); err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Remove - удаляет файл из хранилища.
func (fp *FileProvider) Remove(ctx context.Context, filePath string) error {
	fp.traceCmd(ctx, "Remove", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return err
	}

	if err := os.Remove(fp.rootDir + filePath); err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Ping - проверяет работоспособность файлового хранилища.
// Создаёт временный файл и сразу удаляет его для проверки прав на запись.
// Использует мьютекс для предотвращения одновременных проверок.
func (fp *FileProvider) Ping(ctx context.Context) error {
	fp.traceCmd(ctx, "Ping", testFile)

	fp.muPing.Lock()
	defer fp.muPing.Unlock()

	if dst, err := os.Create(fp.rootDir + testFile); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fp.wrapError(err)
		}
	} else if err = dst.Close(); err != nil {
		return errors.ErrSystemStorageFailedToClose.Wrap(err, "source_test_file", testFile)
	}

	if err := os.Remove(fp.rootDir + testFile); err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Close - закрывает провайдер и освобождает ресурсы.
// Для FileProvider не требует дополнительных действий, возвращает nil.
func (fp *FileProvider) Close() error {
	return nil
}

func (fp *FileProvider) openFile(_ context.Context, filePath string) (*os.File, error) {
	if err := fp.checkFilePath(filePath); err != nil {
		return nil, err
	}

	return os.Open(fp.rootDir + filePath) //nolint:gosec
}

func (fp *FileProvider) getFileInfo(filePath string, fileInfo os.FileInfo) (modelmedia.FileInfo, error) {
	contentType, err := fp.fs.MimeTypes().ContentTypeByExt(path.Ext(filePath))
	if err != nil {
		return modelmedia.FileInfo{}, err
	}

	if fileInfo.Size() < 0 {
		return modelmedia.FileInfo{}, errors.NewInternalError(
			"file size is negative",
			"file", filePath,
		)
	}

	return modelmedia.FileInfo{
		ContentType: contentType,
		Name:        fileInfo.Name(),
		Path:        filePath,
		Size:        fileInfo.Size(),
		UpdatedAt:   casttype.TimeToPointer(fileInfo.ModTime()),
	}, nil
}

func (fp *FileProvider) checkFilePath(filePath string) error {
	length := len(filePath)

	if length < 3 {
		return ErrInternalInvalidPath.New("path", filePath)
	}

	for i := 1; i < length; i++ {
		if filePath[i-1] == '.' && filePath[i] == '.' {
			return ErrInternalInvalidPath.New("path", filePath)
		}
	}

	return nil
}
