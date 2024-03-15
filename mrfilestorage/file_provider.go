package mrfilestorage

import (
	"context"
	"io"
	"os"
	"path"
	"strings"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	providerName = "FileStorage"
)

type (
	FileProvider struct {
		fs      *FileSystem
		rootDir string
	}
)

func NewFileProvider(fs *FileSystem, rootDir string) *FileProvider {
	return &FileProvider{
		fs:      fs,
		rootDir: strings.TrimRight(rootDir, "/") + "/",
	}
}

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

func (fp *FileProvider) DownloadFile(ctx context.Context, filePath string) (io.ReadCloser, error) {
	fp.traceCmd(ctx, "DownloadContent", filePath)

	fd, err := fp.openFile(ctx, filePath)

	if err != nil {
		return nil, fp.wrapError(err)
	}

	return fd, nil
}

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

func (fp *FileProvider) Remove(ctx context.Context, filePath string) error {
	fp.traceCmd(ctx, "Remove", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return err
	}

	return os.Remove(fp.rootDir + filePath)
}

func (fp *FileProvider) openFile(ctx context.Context, filePath string) (*os.File, error) {
	if err := fp.checkFilePath(filePath); err != nil {
		return nil, err
	}

	return os.Open(fp.rootDir + filePath)
}

func (fp *FileProvider) getFileInfo(filePath string, fileInfo os.FileInfo) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType: mrlib.MimeTypeByFile(filePath),
		Name:        fileInfo.Name(),
		Path:        filePath,
		Size:        fileInfo.Size(),
		UpdatedAt:   mrtype.TimePointer(fileInfo.ModTime()),
	}
}

func (fp *FileProvider) checkFilePath(filePath string) error {
	length := len(filePath)

	if length < 3 {
		return FactoryErrInvalidPath.New(filePath)
	}

	for i := 1; i < length; i++ {
		if filePath[i-1] == '.' && filePath[i] == '.' {
			return FactoryErrInvalidPath.New(filePath)
		}
	}

	return nil
}
