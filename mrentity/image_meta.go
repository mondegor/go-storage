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

// Empty - проверяет, что объект пустой.
func (e *ImageMeta) Empty() bool {
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
			return mrcore.ErrInternalTypeAssertion.Wrap(err, "ImageMeta", value)
		}

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("ImageMeta", value)
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
func ImageMetaToInfo(meta ImageMeta, mime *mrlib.MimeTypeList) mrtype.ImageInfo {
	if meta.ContentType == "" && mime != nil {
		meta.ContentType = mime.ContentTypeByFileName(meta.Path)
	}

	return mrtype.ImageInfo{
		ContentType: meta.ContentType,
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Width:     meta.Width,
		Height:    meta.Height,
		Size:      meta.Size,
		CreatedAt: meta.CreatedAt,
		UpdatedAt: meta.UpdatedAt,
	}
}

// ImageMetaToInfoPointer - аналог ImageMetaToInfo, но принимает и возвращает указатель.
func ImageMetaToInfoPointer(meta *ImageMeta, mime *mrlib.MimeTypeList) *mrtype.ImageInfo {
	if meta == nil {
		return nil
	}

	c := ImageMetaToInfo(*meta, mime)

	return &c
}
