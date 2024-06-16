package mrentity

import (
	"database/sql/driver"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ZeronullInt32 - целочисленный тип для которого значение 0 в БД хранится как NULL.
	ZeronullInt32 int32
)

// Scan implements the Scanner interface.
func (n *ZeronullInt32) Scan(value any) error {
	if value == nil {
		*n = 0

		return nil
	}

	if val, ok := value.(int64); ok {
		*n = ZeronullInt32(val)

		return nil
	}

	if val, ok := value.(int32); ok {
		*n = ZeronullInt32(val)

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("ZeronullInt32", value)
}

// Value implements the driver.Valuer interface.
func (n ZeronullInt32) Value() (driver.Value, error) {
	if n == 0 {
		return nil, nil //nolint:nilnil
	}

	return int64(n), nil
}
