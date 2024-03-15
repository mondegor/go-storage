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
	FileMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Size         int64      `json:"size,omitempty"`
		CreatedAt    *time.Time `json:"created,omitempty"`
		UpdatedAt    *time.Time `json:"updated,omitempty"`
	}
)

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
			return mrcore.FactoryErrInternalTypeAssertion.Wrap(err, "FileMeta", value)
		}

		return nil
	}

	return mrcore.FactoryErrInternalTypeAssertion.New("FileMeta", value)
}

// Value implements the driver Valuer interface.
func (n FileMeta) Value() (driver.Value, error) {
	if n.Empty() {
		return nil, nil
	}

	return json.Marshal(n)
}

func FileMetaToInfo(meta FileMeta) mrtype.FileInfo {
	return mrtype.FileInfo{
		ContentType: mrlib.MimeType(meta.ContentType, meta.Path),
		// OriginalName: meta.OriginalName,
		// Name:         path.Base(meta.Path),
		Path:      meta.Path,
		Size:      meta.Size,
		CreatedAt: mrtype.TimePointerCopy(meta.CreatedAt),
		UpdatedAt: mrtype.TimePointerCopy(meta.UpdatedAt),
	}
}

func FileMetaToInfoPointer(meta *FileMeta) *mrtype.FileInfo {
	if meta == nil {
		return nil
	}

	c := FileMetaToInfo(*meta)

	return &c
}
