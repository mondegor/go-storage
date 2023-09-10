package mrentity

import (
    "database/sql/driver"

    "github.com/mondegor/go-sysmess/mrerr"
)

type ZeronullInt64 int64

// Value implements the driver Valuer interface.
func (n ZeronullInt64) Value() (driver.Value, error) {
    if n == 0 {
        return nil, nil
    }

    return int64(n), nil
}

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

    return mrerr.ErrFactoryInternalTypeAssertion.New("ZeronullInt64", value)
}
