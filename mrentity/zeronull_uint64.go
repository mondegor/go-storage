package mrentity

import (
	"database/sql/driver"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// ZeronullUint64 - целочисленный тип для которого значение 0 в БД хранится как NULL.
	ZeronullUint64 uint64
)

// Scan implements the Scanner interface.
func (e *ZeronullUint64) Scan(value any) error {
	if value == nil {
		*e = 0

		return nil
	}

	if val, ok := value.(uint64); ok {
		*e = ZeronullUint64(val)

		return nil
	}

	if val, ok := value.(uint32); ok {
		*e = ZeronullUint64(val)

		return nil
	}

	if val, ok := value.(int64); ok {
		if val < 0 {
			return mrcore.ErrInternalInvalidType.New("int64 < 0", "int64 >= 0")
		}

		*e = ZeronullUint64(val)

		return nil
	}

	if val, ok := value.(int32); ok {
		if val < 0 {
			return mrcore.ErrInternalInvalidType.New("int32 < 0", "int32 >= 0")
		}

		*e = ZeronullUint64(val)

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("ZeronullUint64", value)
}

// Value implements the driver.Valuer interface.
func (e ZeronullUint64) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil //nolint:nilnil
	}

	return uint64(e), nil
}
