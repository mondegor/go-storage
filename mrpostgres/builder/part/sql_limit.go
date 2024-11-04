package part

import (
	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLLimitBuilder - объект для создания части SQL используемой в OFFSET, LIMIT.
	SQLLimitBuilder struct {
		maxSize uint64
		helper  *helper.SQLLimit
	}
)

// NewSQLLimitBuilder - создаёт объект SQLLimitBuilder.
func NewSQLLimitBuilder(maxSize uint64) *SQLLimitBuilder {
	return &SQLLimitBuilder{
		maxSize: maxSize,
		helper:  helper.NewSQLLimit(),
	}
}

// Build - создаёт часть SQL, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLLimitBuilder) Build(index, size uint64) mrstorage.SQLPart {
	return b.createPart(b.helper.OffsetLimit(index, size, b.maxSize))
}

func (b *SQLLimitBuilder) createPart(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return mrsql.NewPart(startArgumentNumber, part)
}
