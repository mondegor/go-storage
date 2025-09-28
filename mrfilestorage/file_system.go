package mrfilestorage

import (
	"errors"
	"fmt"
	"os"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/extfile"
)

type (
	// FileSystem - объект для работы с файлами проекта.
	FileSystem struct {
		dirMode    os.FileMode
		createDirs bool // if not exists
		mimeTypes  *extfile.MimeTypeList
	}
)

// New - создаёт объект FileSystem.
func New(dirMode os.FileMode, createDirs bool, mimeTypes *extfile.MimeTypeList) *FileSystem {
	return &FileSystem{
		dirMode:    dirMode,
		mimeTypes:  mimeTypes,
		createDirs: createDirs,
	}
}

// InitRootDir - comment method.
func (f *FileSystem) InitRootDir(path string) (bool, error) {
	_, err := os.Stat(path)

	if !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return false, mr.ErrInternal.Wrap(err)
		}

		return false, nil
	}

	if !f.createDirs {
		return false, mr.ErrInternal.Wrap(fmt.Errorf("root dir '%s' not exists", path))
	}

	if err = os.Mkdir(path, f.dirMode); err != nil {
		return false, mr.ErrInternal.Wrap(err)
	}

	return true, nil
}

// CreateDirIfNotExists - comment method.
func (f *FileSystem) CreateDirIfNotExists(rootDir, dirPath string) error {
	if _, err := os.Stat(rootDir); err != nil {
		return mr.ErrInternal.Wrap(err)
	}

	dirPath = strings.TrimRight(rootDir, "/") + "/" + strings.Trim(dirPath, "/")

	if _, err := os.Stat(dirPath); !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return mr.ErrInternal.Wrap(err)
		}

		return nil
	}

	return os.MkdirAll(dirPath, f.dirMode)
}

// MimeTypes - comment method.
func (f *FileSystem) MimeTypes() *extfile.MimeTypeList {
	return f.mimeTypes
}
