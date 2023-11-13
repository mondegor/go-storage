package mrpostgres

import (
	"fmt"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	SqlBuilderPager struct {
		maxSize uint64
	}
)

func NewSqlBuilderPager(maxSize uint64) *SqlBuilderPager {
	return &SqlBuilderPager{
		maxSize: maxSize,
	}
}

func (b *SqlBuilderPager) OffsetLimit(index, size uint64) mrstorage.SqlBuilderPartFunc {
	if b.maxSize > 0 && (size < 1 || size > b.maxSize) {
		size = b.maxSize
	} else if size < 1 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		if index > 0 {
			return fmt.Sprintf(
				" OFFSET %d LIMIT %d",
				index * size,
				size,
			), []any{}
		}

		return fmt.Sprintf(" LIMIT %d", size), []any{}
	}
}
