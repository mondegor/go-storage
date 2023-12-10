package mrfilestorage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	FileProvider struct {
		fs      *FileSystem
		rootDir string
		baseDir string
	}
)

func NewFileProvider(fs *FileSystem, rootDir string) *FileProvider {
	path := strings.TrimRight(rootDir, "/") + "/"

	return &FileProvider{
		fs:      fs,
		rootDir: path,
		baseDir: path,
	}
}

func (fp *FileProvider) WithBaseDir(value string) (mrstorage.ExtFileProviderAPI, error) {
	if value != "" {
		value = strings.Trim(value, "/")

		if value != "" {
			value += "/"
		}
	}

	if err := fp.fs.CreateDirIfNotExists(fp.rootDir, value); err != nil {
		return nil, err
	}

	c := *fp
	c.baseDir = fp.rootDir + value

	return &c, nil
}

func (fp *FileProvider) Info(ctx context.Context, fileName string) (mrtype.FileInfo, error) {
	return fp.getFileInfo(fp.baseDir + fileName)
}

func (fp *FileProvider) Download(ctx context.Context, fileName string) (*mrtype.File, error) {
	fileInfo, err := fp.getFileInfo(fp.baseDir + fileName)

	if err != nil {
		return nil, err
	}

	fd, err := os.Open(fp.baseDir + fileName)

	if err != nil {
		return nil, fp.wrapError(err, 0)
	}

	return &mrtype.File{
		FileInfo: fileInfo,
		Path:     fileName,
		Body:     fd,
	}, nil
}

func (fp *FileProvider) Upload(ctx context.Context, file *mrtype.File) error {
	dst, err := os.Create(fp.baseDir + file.Path)

	if err != nil {
		return fp.wrapError(err, 0)
	}

	defer dst.Close()

	if _, err = io.Copy(dst, file.Body); err != nil {
		return fp.wrapError(err, 0)
	}

	return nil
}

func (fp *FileProvider) Remove(ctx context.Context, fileName string) error {
	return os.Remove(fp.baseDir + fileName)
}

func (fp *FileProvider) getFileInfo(filePath string) (mrtype.FileInfo, error) {
	fi, err := os.Stat(filePath)

	if err != nil {
		return mrtype.FileInfo{}, fp.wrapError(err, 1)
	}

	return mrtype.FileInfo{
		ContentType:  mrlib.MimeTypeByExt(filepath.Ext(filePath)),
		Name:         filepath.Base(filePath),
		LastModified: fi.ModTime(),
		Size:         fi.Size(),
	}, nil
}
