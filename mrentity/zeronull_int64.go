package mrentity

import (
	"database/sql/driver"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ZeronullInt64 - comment type.
	ZeronullInt64 int64
)

// Scan implements the Scanner interface.
func (n *ZeronullInt64) Scan(value any) error {
	if value == nil {
		*n = 0

		return nil
	}

	if val, ok := value.(int64); ok {
		*n = ZeronullInt64(val)

		return nil
	}

	if val, ok := value.(int32); ok {
		*n = ZeronullInt64(val)

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("ZeronullInt64", value)
}

// Value implements the driver.Valuer interface.
func (n ZeronullInt64) Value() (driver.Value, error) {
	if n == 0 {
		return nil, nil //nolint:nilnil
	}

	return int64(n), nil
}
