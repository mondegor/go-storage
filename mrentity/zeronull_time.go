package mrentity

import (
    "database/sql/driver"
    "time"

    "github.com/mondegor/go-webcore/mrcore"
)

type (
	ZeronullTime time.Time
)

// Value implements the driver Valuer interface.
func (n ZeronullTime) Value() (driver.Value, error) {
    if time.Time(n).IsZero() {
        return nil, nil
    }

    return time.Time(n), nil
}

// Scan implements the Scanner interface.
func (n *ZeronullTime) Scan(value any) error {
    if value == nil {
        *n = ZeronullTime{}
        return nil
    }

    if val, ok := value.(time.Time); ok {
        *n = ZeronullTime(val)
        return nil
    }

    return mrcore.FactoryErrInternalTypeAssertion.New("ZeronullTime", value)
}
