package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlib/copyptr"
	"github.com/mondegor/go-sysmess/mrtype"
)

type (
	// ImageMeta - метаинформация об изображении, позволяет сохранять в БД и читать из неё в виде json.
	ImageMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Width        uint64     `json:"width,omitempty"`
		Height       uint64     `json:"height,omitempty"`
		Size         uint64     `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		UpdatedAt    *time.Time `json:"updated,omitempty"`
	}
)

// Empty - сообщает, является ли объект пустым.
func (e ImageMeta) Empty() bool {
	return e.Path == "" &&
		e.OriginalName == ""
}

// Scan implements the Scanner interface.
func (e *ImageMeta) Scan(value any) error {
	if value == nil {
		*e = ImageMeta{}

		return nil
	}

	if val, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(val), e); err != nil {
			return mr.ErrInternalTypeAssertion.Wrap(err, "ImageMeta", value)
		}

		return nil
	}

	return mr.ErrInternalTypeAssertion.New("ImageMeta", value)
}

// Value implements the driver.Valuer interface.
func (e ImageMeta) Value() (driver.Value, error) {
	if e.Empty() {
		return nil, nil //nolint:nilnil
	}

	return json.Marshal(e)
}

// ImageMetaToInfo - преобразование данных изображения предназначенных
// для хранилища в формат данных для передачи клиенту.
func ImageMetaToInfo(meta ImageMeta) mrtype.ImageInfo {
	return mrtype.ImageInfo{
		ContentType: meta.ContentType,
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Width:     meta.Width,
		Height:    meta.Height,
		Size:      meta.Size,
		CreatedAt: copyptr.Time(meta.CreatedAt),
		UpdatedAt: copyptr.Time(meta.UpdatedAt),
	}
}

// ImageMetaToInfoPointer - аналог ImageMetaToInfo, но принимает и возвращает указатель.
func ImageMetaToInfoPointer(meta *ImageMeta) *mrtype.ImageInfo {
	if meta == nil {
		return nil
	}

	c := ImageMetaToInfo(*meta)

	return &c
}
