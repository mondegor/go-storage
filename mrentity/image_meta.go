package mrentity

import (
	"database/sql/driver"
	"encoding/json"
	"path"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	ImageMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Width        int32      `json:"w,omitempty"`
		Height       int32      `json:"h,omitempty"`
		Size         int64      `json:"s,omitempty"`
		CreatedAt    *time.Time `json:"crt,omitempty"`
		UpdatedAt    *time.Time `json:"upd,omitempty"`
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

func ConvertImageMetaToInfo(meta *ImageMeta) *mrtype.ImageInfo {
	if meta == nil {
		return nil
	}

	return &mrtype.ImageInfo{
		ContentType:  meta.ContentType,
		OriginalName: meta.OriginalName,
		Name:         path.Base(meta.Path),
		Path:         meta.Path,
		Width:        meta.Width,
		Height:       meta.Height,
		Size:         meta.Size,
		CreatedAt:    meta.CreatedAt,
		ModifiedAt:   meta.UpdatedAt,
	}
}
