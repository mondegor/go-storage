package mrfilestorage

import (
	"context"
	"errors"
	"io"
	"os"
	"path"
	"strings"
	"sync"

	"github.com/mondegor/go-webcore/mrtype"
)

const (
	providerName = "FileStorage"
	testFile     = "testFile-d6b6943c-e1b2-4625-b133-9805a5cf5f8d"
)

type (
	// FileProvider - файловый провайдер, работающий с нативной файловой системой,
	// позволяет читать, сохранять, удалять файлы.
	FileProvider struct {
		fs      *FileSystem
		rootDir string
		muPing  sync.Mutex
	}
)

// NewFileProvider - создаёт объект FileProvider.
func NewFileProvider(fs *FileSystem, rootDir string) *FileProvider {
	return &FileProvider{
		fs:      fs,
		rootDir: strings.TrimRight(rootDir, "/") + "/",
	}
}

// Info - comment method.
func (fp *FileProvider) Info(ctx context.Context, filePath string) (mrtype.FileInfo, error) {
	fp.traceCmd(ctx, "Info", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return mrtype.FileInfo{}, err
	}

	fi, err := os.Stat(fp.rootDir + filePath)
	if err != nil {
		return mrtype.FileInfo{}, fp.wrapError(err)
	}

	return fp.getFileInfo(filePath, fi), nil
}

// Download - comment method.
func (fp *FileProvider) Download(ctx context.Context, filePath string) (mrtype.File, error) {
	fp.traceCmd(ctx, "Download", filePath)

	fd, err := fp.openFile(ctx, filePath)
	if err != nil {
		return mrtype.File{}, fp.wrapError(err)
	}

	fi, err := fd.Stat()
	if err != nil {
		return mrtype.File{}, fp.wrapError(err)
	}

	return mrtype.File{
		FileInfo: fp.getFileInfo(filePath, fi),
		Body:     fd,
	}, nil
}

// DownloadFile - comment method.
func (fp *FileProvider) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fp.traceCmd(ctx, "DownloadContent", filePath)

	fd, err := fp.openFile(ctx, filePath)
	if err != nil {
		return nil, fp.wrapError(err)
	}

	return fd, nil
}

// Upload - comment method.
func (fp *FileProvider) Upload(ctx context.Context, file mrtype.File) error {
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

	defer dst.Close()

	if _, err = io.Copy(dst, file.Body); err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Remove - comment method.
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

// Ping - проверяет возможность работы с файлами.
func (fp *FileProvider) Ping(ctx context.Context) error {
	fp.traceCmd(ctx, "Ping", testFile)

	fp.muPing.Lock()
	defer fp.muPing.Unlock()

	if dst, err := os.Create(fp.rootDir + testFile); err != nil {
		if !errors.Is(err, os.ErrExist) {
			return fp.wrapError(err)
		}
	} else if err = dst.Close(); err != nil {
		return fp.wrapError(err)
	}

	if err := os.Remove(fp.rootDir + testFile); err != nil {
		return fp.wrapError(err)
	}

	return nil
}

// Close - закрывает текущее соединение.
func (fp *FileProvider) Close() error {
	return nil
}

func (fp *FileProvider) openFile(_ context.Context, filePath string) (*os.File, error) {
	if err := fp.checkFilePath(filePath); err != nil {
		return nil, err
	}

	return os.Open(fp.rootDir + filePath)
}

func (fp *FileProvider) getFileInfo(filePath string, fileInfo os.FileInfo) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType: fp.fs.MimeTypes().ContentTypeByFileName(filePath),
		Name:        fileInfo.Name(),
		Path:        filePath,
		Size:        fileInfo.Size(),
		UpdatedAt:   mrtype.TimeToPointer(fileInfo.ModTime()),
	}
}

func (fp *FileProvider) checkFilePath(filePath string) error {
	length := len(filePath)

	if length < 3 {
		return ErrInvalidPath.New(filePath)
	}

	for i := 1; i < length; i++ {
		if filePath[i-1] == '.' && filePath[i] == '.' {
			return ErrInvalidPath.New(filePath)
		}
	}

	return nil
}
