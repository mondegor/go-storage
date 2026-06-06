package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"time"

	"github.com/mondegor/go-sysmess/errors"
	modelmedia "github.com/mondegor/go-sysmess/mrmodel/media"
	"github.com/mondegor/go-sysmess/util/copyptr"
)

type (
	// ImageMeta - метаинформация об изображении, позволяет сохранять в БД и читать из неё в виде json.
	// Реализует интерфейсы sql.Scanner и driver.Valuer.
	ImageMeta struct {
		Path         string     `json:"path,omitempty"`    // Path - путь к изображению в хранилище
		ContentType  string     `json:"type,omitempty"`    // ContentType - MIME-тип изображения (например: "image/png")
		OriginalName string     `json:"origin,omitempty"`  // OriginalName - оригинальное имя файла при загрузке
		Width        int32      `json:"width,omitempty"`   // Width - ширина изображения в пикселях
		Height       int32      `json:"height,omitempty"`  // Height - высота изображения в пикселях
		Size         int64      `json:"size,omitempty"`    // Size - размер файла в байтах
		CreatedAt    *time.Time `json:"created,omitempty"` // CreatedAt - время создания записи
		UpdatedAt    *time.Time `json:"updated,omitempty"` // UpdatedAt - время последнего обновления записи
	}
)

// Empty - сообщает, является ли объект пустым.
// Проверяет только Path и OriginalName, так как эти поля являются обязательными идентификаторами изображения.
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
			return errors.ErrInternalTypeAssertion.Wrap(
				err,
				"type", "ImageMeta",
				"value", value,
			)
		}

		return nil
	}

	return errors.ErrInternalTypeAssertion.New(
		"type", "ImageMeta",
		"value", value,
	)
}

// Value implements the driver.Valuer interface.
func (e ImageMeta) Value() (driver.Value, error) {
	if e.Empty() {
		return nil, nil //nolint:nilnil
	}

	return json.Marshal(e)
}

// ImageMetaToInfo - преобразование данных изображения, предназначенных для хранения в БД,
// в формат данных для передачи клиенту (mrmodel.ImageInfo).
func ImageMetaToInfo(meta ImageMeta) modelmedia.ImageInfo {
	return modelmedia.ImageInfo{
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

// ImageMetaToInfoPointer - аналог ImageMetaToInfo, но работает с указателями.
func ImageMetaToInfoPointer(meta *ImageMeta) *modelmedia.ImageInfo {
	if meta == nil {
		return nil
	}

	c := ImageMetaToInfo(*meta)

	return &c
}
