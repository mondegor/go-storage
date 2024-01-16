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
	fp.debugCmd(ctx, "Info", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return mrtype.FileInfo{}, err
	}

	return fp.getFileInfo(filePath)
}

func (fp *FileProvider) Download(ctx context.Context, filePath string) (mrtype.File, error) {
	fp.debugCmd(ctx, "Download", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return mrtype.File{}, err
	}

	fileInfo, err := fp.getFileInfo(filePath)

	if err != nil {
		return mrtype.File{}, err
	}

	fd, err := os.Open(fp.rootDir + filePath)

	if err != nil {
		return mrtype.File{}, fp.wrapError(err, 0)
	}

	return mrtype.File{
		FileInfo: fileInfo,
		Body:     fd,
	}, nil
}

func (fp *FileProvider) Upload(ctx context.Context, file mrtype.File) error {
	fp.debugCmd(ctx, "Upload", file.Path)

	if err := fp.checkFilePath(file.Path); err != nil {
		return err
	}

	if dirPath := path.Dir(file.Path); dirPath != "" {
		if err := fp.fs.CreateDirIfNotExists(fp.rootDir, dirPath); err != nil {
			return fp.wrapError(err, 0)
		}
	}

	dst, err := os.Create(fp.rootDir + file.Path)

	if err != nil {
		return fp.wrapError(err, 0)
	}

	defer dst.Close()

	if _, err = io.Copy(dst, file.Body); err != nil {
		return fp.wrapError(err, 0)
	}

	return nil
}

func (fp *FileProvider) Remove(ctx context.Context, filePath string) error {
	fp.debugCmd(ctx, "Remove", filePath)

	if err := fp.checkFilePath(filePath); err != nil {
		return err
	}

	return os.Remove(fp.rootDir + filePath)
}

func (fp *FileProvider) getFileInfo(filePath string) (mrtype.FileInfo, error) {
	fi, err := os.Stat(fp.rootDir + filePath)

	if err != nil {
		return mrtype.FileInfo{}, fp.wrapError(err, 1)
	}

	return mrtype.FileInfo{
		ContentType: mrlib.MimeTypeByFile(filePath),
		Name:        path.Base(filePath),
		Path:        filePath,
		Size:        fi.Size(),
		ModifiedAt:  mrtype.TimePointer(fi.ModTime()),
	}, nil
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
