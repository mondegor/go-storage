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
	// FileMeta - comment struct.
	FileMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		UpdatedAt    *time.Time `json:"updated,omitempty"`
	}
)

// Empty - проверяет, что объект пустой.
func (n *FileMeta) Empty() bool {
	return n.Path == "" &&
		n.OriginalName == ""
}

// Scan implements the Scanner interface.
func (n *FileMeta) Scan(value any) error {
	if value == nil {
		*n = FileMeta{}

		return nil
	}

	if val, ok := value.(string); ok {
		if err := json.Unmarshal([]byte(val), n); err != nil {
			return mrcore.ErrInternalTypeAssertion.Wrap(err, "FileMeta", value)
		}

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("FileMeta", value)
}

// Value implements the driver.Valuer interface.
func (n FileMeta) Value() (driver.Value, error) {
	if n.Empty() {
		return nil, nil //nolint:nilnil
	}

	return json.Marshal(n)
}

// FileMetaToInfo - comment func.
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
		CreatedAt: mrtype.TimePointerCopy(meta.CreatedAt),
		UpdatedAt: mrtype.TimePointerCopy(meta.UpdatedAt),
	}
}

// FileMetaToInfoPointer - comment func.
func FileMetaToInfoPointer(meta *FileMeta, mime *mrlib.MimeTypeList) *mrtype.FileInfo {
	if meta == nil {
		return nil
	}

	c := FileMetaToInfo(*meta, mime)

	return &c
}
