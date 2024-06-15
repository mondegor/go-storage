package mrentity

import (
	"database/sql/driver"

	"github.com/mondegor/go-webcore/mrcore"
)

type (
	// EmptynullString - comment type.
	EmptynullString string
)

// Scan implements the Scanner interface.
func (n *EmptynullString) Scan(value any) error {
	if value == nil {
		*n = ""

		return nil
	}

	if val, ok := value.(string); ok {
		*n = EmptynullString(val)

		return nil
	}

	return mrcore.ErrInternalTypeAssertion.New("EmptynullString", value)
}

// Value implements the driver.Valuer interface.
func (n EmptynullString) Value() (driver.Value, error) {
	if n == "" {
		return nil, nil //nolint:nilnil
	}

	return string(n), nil
}
