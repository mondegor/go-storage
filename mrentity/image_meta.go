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
	ImageMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Width        int32      `json:"width,omitempty"`
		Height       int32      `json:"height,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		UpdatedAt    *time.Time `json:"updated,omitempty"`
	}
)

func (n *ImageMeta) Empty() bool {
	return n.Path == "" &&
		n.OriginalName == ""
}

// Scan implements the Scanner interface.
func (n *ImageMeta) Scan(value any) error {
	if value == nil {
		*n = ImageMeta{}
		return nil
	}

	if val, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(val), n); err != nil {
			return mrcore.FactoryErrInternalTypeAssertion.Wrap(err, "ImageMeta", value)
		}

		return nil
	}

	return mrcore.FactoryErrInternalTypeAssertion.New("ImageMeta", value)
}

// Value implements the driver Valuer interface.
func (n ImageMeta) Value() (driver.Value, error) {
	if n.Empty() {
		return nil, nil
	}

	return json.Marshal(n)
}

func ImageMetaToInfo(meta ImageMeta) mrtype.ImageInfo {
	return mrtype.ImageInfo{
		ContentType: mrlib.MimeType(meta.ContentType, meta.Path),
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

func ImageMetaToInfoPointer(meta *ImageMeta) *mrtype.ImageInfo {
	if meta == nil {
		return nil
	}

	c := ImageMetaToInfo(*meta)

	return &c
}
