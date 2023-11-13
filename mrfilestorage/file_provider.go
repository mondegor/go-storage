package mrfilestorage

import (
	"context"
	"io"
	"os"
	"path/filepath"
	"strings"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	nativeAdapter struct {
		rootDir string
	}
)

func New(rootDir string) *nativeAdapter {
	return &nativeAdapter{
		rootDir: strings.TrimRight(rootDir, "/"),
	}
}

func (n *nativeAdapter) Info(ctx context.Context, path string) (mrtype.FileInfo, error) {
	return n.getFileInfo(n.getFullPath(path))
}

func (n *nativeAdapter) Download(ctx context.Context, path string) (*mrtype.File, error) {
	fullPath := n.getFullPath(path)
	fileInfo, err := n.getFileInfo(fullPath)

	if err != nil {
		return nil, err
	}

	fd, err := os.Open(fullPath)

	if err != nil {
		return nil, err
	}

	return &mrtype.File{
		FileInfo: fileInfo,
		Path:     path,
		Body:     fd,
	}, nil
}

func (n *nativeAdapter) Upload(ctx context.Context, file *mrtype.File) error {
	dst, err := os.Create(n.getFullPath(file.Path))

	if err != nil {
		return err
	}

	defer dst.Close()

	_, err = io.Copy(dst, file.Body)

	return err
}

func (n *nativeAdapter) Remove(ctx context.Context, path string) error {
	return os.Remove(n.getFullPath(path))
}

func (n *nativeAdapter) getFullPath(path string) string {
	return strings.Join([]string{n.rootDir, path}, "/")
}

func (n *nativeAdapter) getFileInfo(path string) (mrtype.FileInfo, error) {
	fi, err := os.Stat(path)

	if err != nil {
		return mrtype.FileInfo{}, err
	}

	return mrtype.FileInfo{
		ContentType:  mrlib.MimeTypeByExt(filepath.Ext(path)),
		Name:         filepath.Base(path),
		LastModified: fi.ModTime(),
		Size:         fi.Size(),
	}, nil
}
