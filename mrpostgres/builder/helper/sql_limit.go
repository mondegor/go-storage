package helper

import (
	"strconv"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLLimit - объект для создания независимой части SQL используемой при создании лимита (OFFSET, LIMIT).
	SQLLimit struct{}
)

// NewSQLLimit - создаёт объект SQLLimit.
func NewSQLLimit() *SQLLimit {
	return &SQLLimit{}
}

// OffsetLimit - возвращает SQL лимит с указанными значениями.
// При size = 0 лимит или ограничен maxSize или не ограничен, если maxSize = 0.
func (b *SQLLimit) OffsetLimit(index, size, maxSize uint64) mrstorage.SQLPartFunc {
	if maxSize > 0 && (size == 0 || size > maxSize) {
		size = maxSize
	} else if size == 0 {
		return nil
	}

	return func(_ int) (string, []any) {
		if index > 0 {
			return " OFFSET " + strconv.FormatUint(index*size, 10) +
				" LIMIT " + strconv.FormatUint(size, 10), nil
		}

		return " LIMIT " + strconv.FormatUint(size, 10), nil
	}
}
