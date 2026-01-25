package mrfilestorage

import (
	"os"
	"strings"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/util/mime"
)

type (
	// FileSystem - объект для работы с файлами проекта.
	FileSystem struct {
		dirMode    os.FileMode
		createDirs bool // if not exists
		mimeTypes  *mime.TypeList
	}
)

// New - создаёт объект FileSystem.
func New(dirMode os.FileMode, createDirs bool, mimeTypes *mime.TypeList) *FileSystem {
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
			return false, errors.WrapInternalError(err, "InitRootDir:os.Stat failed")
		}

		return false, nil
	}

	if !f.createDirs {
		return false, errors.NewInternalError("root dir not exists", "dir", path)
	}

	if err = os.Mkdir(path, f.dirMode); err != nil {
		return false, errors.WrapInternalError(err, "InitRootDir:os.Mkdir failed")
	}

	return true, nil
}

// CreateDirIfNotExists - comment method.
func (f *FileSystem) CreateDirIfNotExists(rootDir, dirPath string) error {
	if _, err := os.Stat(rootDir); err != nil {
		return errors.WrapInternalError(err, "CreateDirIfNotExists:os.Stat rootDir failed")
	}

	dirPath = strings.TrimRight(rootDir, "/") + "/" + strings.Trim(dirPath, "/")

	if _, err := os.Stat(dirPath); !errors.Is(err, os.ErrNotExist) {
		if err != nil {
			return errors.WrapInternalError(err, "CreateDirIfNotExists:os.Stat dirPath failed")
		}

		return nil
	}

	return os.MkdirAll(dirPath, f.dirMode)
}

// MimeTypes - comment method.
func (f *FileSystem) MimeTypes() *mime.TypeList {
	return f.mimeTypes
}
