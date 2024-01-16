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
	FileMeta struct {
		Path         string     `json:"path,omitempty"`
		ContentType  string     `json:"type,omitempty"`
		OriginalName string     `json:"origin,omitempty"`
		Size         int64      `json:"s,omitempty"`
		CreatedAt    *time.Time `json:"crt,omitempty"`
		UpdatedAt    *time.Time `json:"upd,omitempty"`
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

func ConvertFileMetaToInfo(meta *FileMeta) *mrtype.FileInfo {
	if meta == nil {
		return nil
	}

	return &mrtype.FileInfo{
		ContentType:  meta.ContentType,
		OriginalName: meta.OriginalName,
		Name:         path.Base(meta.Path),
		Path:         meta.Path,
		Size:         meta.Size,
		CreatedAt:    meta.CreatedAt,
		ModifiedAt:   meta.UpdatedAt,
	}
}
