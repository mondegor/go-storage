package mrpostgres

import (
	"fmt"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderLimit - comment struct.
	SQLBuilderLimit struct {
		maxSize uint64
	}
)

// NewSQLBuilderLimit - создаёт объект SQLBuilderLimit.
func NewSQLBuilderLimit(maxSize uint64) *SQLBuilderLimit {
	return &SQLBuilderLimit{
		maxSize: maxSize,
	}
}

// OffsetLimit - comment method.
func (b *SQLBuilderLimit) OffsetLimit(index, size uint64) mrstorage.SQLBuilderPartFunc {
	if b.maxSize > 0 && (size < 1 || size > b.maxSize) {
		size = b.maxSize
	} else if size < 1 {
		return nil
	}

	return func(_ int) (string, []any) {
		if index > 0 {
			return fmt.Sprintf(
				" OFFSET %d LIMIT %d",
				index*size,
				size,
			), nil
		}

		return fmt.Sprintf(" LIMIT %d", size), nil
	}
}
