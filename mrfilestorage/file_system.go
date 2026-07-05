package mrfilestorage

import (
	"os"
	"strings"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/util/mime"
)

type (
	// FileSystem - объект для работы с файлами проекта.
	// Предоставляет утилитарные методы для управления директориями и работы с MIME-типами.
	FileSystem struct {
		dirMode    os.FileMode    // dirMode - режим доступа для создаваемых директорий (например: 0755)
		createDirs bool           // createDirs - флаг автоматического создания директорий, если они не существуют
		mimeTypes  *mime.TypeList // mimeTypes - список MIME-типов для определения типа контента по расширению
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

// InitRootDir - инициализирует корневую директорию.
// Если директория не существует и createDirs=true, создаёт её.
// Возвращает true, если директория была создана, false - если уже существовала.
func (f *FileSystem) InitRootDir(path string) (bool, error) {
	_, err := os.Stat(path)
	if err == nil {
		return false, nil // директория уже существует
	}

	if !errors.Is(err, os.ErrNotExist) {
		return false, errors.WrapInternalError(err, "InitRootDir:os.Stat failed")
	}

	// здесь err == os.ErrNotExist

	if !f.createDirs {
		return false, errors.NewInternalError("root dir not exists", "dir", path)
	}

	if err = os.Mkdir(path, f.dirMode); err != nil {
		return false, errors.WrapInternalError(err, "InitRootDir:os.Mkdir failed")
	}

	return true, nil // директория была создана
}

// CreateDirIfNotExists - создаёт директорию по указанному пути, если она не существует.
// Формирует полный путь из rootDir и dirPath. Проверяет существование rootDir перед созданием.
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

// MimeTypes - возвращает список MIME-типов, с которыми должна работать файловая система.
func (f *FileSystem) MimeTypes() *mime.TypeList {
	return f.mimeTypes
}
