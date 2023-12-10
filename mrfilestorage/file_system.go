package mrfilestorage

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	FileSystem struct {
		dirMode    os.FileMode
		createDirs bool // if not exists
	}

	Options struct {
		DirMode    os.FileMode
		CreateDirs bool
	}
)

func New(opt Options) *FileSystem {
	return &FileSystem{
		dirMode:    opt.DirMode,
		createDirs: opt.CreateDirs,
	}
}

func (f *FileSystem) InitRootDir(path string) (bool, error) {
	_, err := os.Stat(path)

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return false, mrcore.FactoryErrInternal.Wrap(err)
		}

		return false, nil
	}

	if !f.createDirs {
		return false, fmt.Errorf("root dir '%s' not exists", path)
	}

	if err = os.Mkdir(path, f.dirMode); err != nil {
		return false, mrcore.FactoryErrInternal.Wrap(err)
	}

	return true, nil
}

func (f *FileSystem) CreateDirIfNotExists(rootDir string, dirPath string) error {
	if _, err := os.Stat(rootDir); err != nil {
		return mrcore.FactoryErrInternal.Wrap(err)
	}

	dirPath = strings.TrimRight(rootDir, "/") + "/" + strings.Trim(dirPath, "/")

	_, err := os.Stat(dirPath)

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return mrcore.FactoryErrInternal.Wrap(err)
		}

		return nil
	}

	return os.MkdirAll(dirPath, f.dirMode)
}
