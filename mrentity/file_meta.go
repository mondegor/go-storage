package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/mondegor/go-sysmess/mrdto"
	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/copyptr"
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

// Empty - сообщает, является ли объект пустым.
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
			return mr.ErrInternalTypeAssertion.Wrap(err, "FileMeta", value)
		}

		return nil
	}

	return mr.ErrInternalTypeAssertion.New("FileMeta", value)
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
func FileMetaToInfo(meta FileMeta) mrdto.FileInfo {
	return mrdto.FileInfo{
		ContentType: meta.ContentType,
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Size:      meta.Size,
		CreatedAt: copyptr.Time(meta.CreatedAt),
		UpdatedAt: copyptr.Time(meta.UpdatedAt),
	}
}

// FileMetaToInfoPointer - аналог FileMetaToInfo, но принимает и возвращает указатель.
func FileMetaToInfoPointer(meta *FileMeta) *mrdto.FileInfo {
	if meta == nil {
		return nil
	}

	c := FileMetaToInfo(*meta)

	return &c
}
