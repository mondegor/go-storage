package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	"github.com/mondegor/go-sysmess/mrmodel"
	"github.com/mondegor/go-sysmess/util/copyptr"
)

type (
	// FileMeta - метаинформация о файле, позволяет сохранять в БД и читать из неё в виде json.
	// Реализует интерфейсы sql.Scanner и driver.Valuer.
	FileMeta struct {
		Path         string     `json:"path,omitempty"`    // Path - путь к файлу в хранилище
		ContentType  string     `json:"type,omitempty"`    // ContentType - MIME-тип файла (например: "image/png")
		OriginalName string     `json:"origin,omitempty"`  // OriginalName - оригинальное имя файла при загрузке
		Size         int64      `json:"size,omitempty"`    // Size - размер файла в байтах
		CreatedAt    *time.Time `json:"created,omitempty"` // CreatedAt - время создания записи
		UpdatedAt    *time.Time `json:"updated,omitempty"` // UpdatedAt - время последнего обновления записи
	}
)

// Empty - сообщает, является ли объект пустым.
// Проверяет только Path и OriginalName, так как эти поля являются обязательными идентификаторами файла.
func (e FileMeta) Empty() bool {
	return e.Path == "" &&
		e.OriginalName == ""
}

// Scan implements the Scanner interface.
func (e *FileMeta) Scan(value any) error {
	if value == nil {
		*e = FileMeta{}

		return nil
	}

	if val, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(val), e); err != nil {
			return errors.ErrInternalTypeAssertion.Wrap(
				err,
				"type", "FileMeta",
				"value", value,
			)
		}

		return nil
	}

	return errors.ErrInternalTypeAssertion.New(
		"type", "FileMeta",
		"value", value,
	)
}

// Value implements the driver.Valuer interface.
func (e FileMeta) Value() (driver.Value, error) {
	if e.Empty() {
		return nil, nil //nolint:nilnil
	}

	return json.Marshal(e)
}

// FileMetaToInfo - преобразование данных файла, предназначенных для хранения в БД,
// в формат данных для передачи клиенту (mrmodel.FileInfo).
func FileMetaToInfo(meta FileMeta) mrmodel.FileInfo {
	return mrmodel.FileInfo{
		ContentType: meta.ContentType,
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Size:      meta.Size,
		CreatedAt: copyptr.Time(meta.CreatedAt),
		UpdatedAt: copyptr.Time(meta.UpdatedAt),
	}
}

// FileMetaToInfoPointer - аналог FileMetaToInfo, но работает с указателями.
func FileMetaToInfoPointer(meta *FileMeta) *mrmodel.FileInfo {
	if meta == nil {
		return nil
	}

	c := FileMetaToInfo(*meta)

	return &c
}
