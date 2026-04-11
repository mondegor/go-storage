package mrentity

import (
	"database/sql/driver"

	"github.com/mondegor/go-sysmess/errors"
)

type (
	// ZeronullUint64 - целочисленный тип для которого значение 0 в БД хранится как NULL.
	// Реализует интерфейсы sql.Scanner и driver.Valuer.
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
			return errors.ErrInternalInvalidType.New(
				"type", "int64 < 0",
				"expected", "int64 >= 0",
			)
		}

		*e = ZeronullUint64(val)

		return nil
	}

	if val, ok := value.(int32); ok {
		if val < 0 {
			return errors.ErrInternalInvalidType.New(
				"type", "int32 < 0",
				"expected", "int32 >= 0",
			)
		}

		*e = ZeronullUint64(val)

		return nil
	}

	return errors.ErrInternalTypeAssertion.New(
		"type", "ZeronullUint64",
		"value", value,
	)
}

// Value implements the driver.Valuer interface.
func (e ZeronullUint64) Value() (driver.Value, error) {
	if e == 0 {
		return nil, nil //nolint:nilnil
	}

	return uint64(e), nil
}
