package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// FileMeta - метаинформация о файле, позволяет сохранять в БД и читать из неё в виде json.
	FileMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Size         uint64     `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		UpdatedAt    *time.Time `json:"updated,omitempty"`
	}
)

// Empty - проверяет, что объект пустой.
func (e *FileMeta) Empty() bool {
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
			return mrcore.ErrInternalTypeAssertion.Wrap(err, "FileMeta", value)
		}

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("FileMeta", value)
}

// Value implements the driver.Valuer interface.
func (e FileMeta) Value() (driver.Value, error) {
	if e.Empty() {
		return nil, nil //nolint:nilnil
	}

	return json.Marshal(e)
}

// FileMetaToInfo - преобразование данных файла предназначенных
// для хранилища в формат данных для передачи клиенту.
func FileMetaToInfo(meta FileMeta, mime *mrlib.MimeTypeList) mrtype.FileInfo {
	if meta.ContentType == "" && mime != nil {
		meta.ContentType = mime.ContentTypeByFileName(meta.Path)
	}

	return mrtype.FileInfo{
		ContentType: meta.ContentType,
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Size:      meta.Size,
		CreatedAt: mrtype.CopyTimePointer(meta.CreatedAt),
		UpdatedAt: mrtype.CopyTimePointer(meta.UpdatedAt),
	}
}

// FileMetaToInfoPointer - аналог FileMetaToInfo, но принимает и возвращает указатель.
func FileMetaToInfoPointer(meta *FileMeta, mime *mrlib.MimeTypeList) *mrtype.FileInfo {
	if meta == nil {
		return nil
	}

	c := FileMetaToInfo(*meta, mime)

	return &c
}
