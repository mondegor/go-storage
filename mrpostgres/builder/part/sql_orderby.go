package part

import (
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLOrderByBuilder - объект для создания части SQL используемой в ORDER BY (field ASC, ...).
	SQLOrderByBuilder struct {
		helper          *helper.SQLOrderBy
		defaultPartFunc mrstorage.SQLPartFunc
	}
)

// NewSQLOrderByBuilder - создаёт объект SQLOrderByBuilder.
func NewSQLOrderByBuilder(defaultSort mrtype.SortParams) *SQLOrderByBuilder {
	b := &SQLOrderByBuilder{
		helper: helper.NewSQLOrderBy(),
	}

	if defaultSort.FieldName != "" {
		b.defaultPartFunc = b.helper.Field(defaultSort.FieldName, defaultSort.Direction)
	}

	return b
}

// Build - создаёт часть SQL, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLOrderByBuilder) Build(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	if part == nil {
		part = b.defaultPartFunc
	}

	return b.createPart(part)
}

// BuildComma - создаёт часть SQL объединяющую независимые части через запятую, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLOrderByBuilder) BuildComma(parts ...mrstorage.SQLPartFunc) mrstorage.SQLPart {
	parts = mrsql.SQLPartFuncRemoveNil(parts)

	if len(parts) == 0 {
		return b.createPart(b.defaultPartFunc)
	}

	if len(parts) == 1 {
		return b.createPart(parts[0])
	}

	return b.createPart(b.helper.JoinComma(parts...))
}

// BuildFunc - создаёт часть SQL с использованием помощника, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLOrderByBuilder) BuildFunc(fn func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc) mrstorage.SQLPart {
	var partFunc mrstorage.SQLPartFunc

	if fn != nil {
		partFunc = fn(b.helper)
	}

	if partFunc != nil {
		return b.createPart(partFunc)
	}

	return b.createPart(b.defaultPartFunc)
}

// HelpFunc - создаёт независимую часть SQL, которая может быть использована при создании других частей SQL.
func (b *SQLOrderByBuilder) HelpFunc(fn func(o mrstorage.SQLOrderByHelper) mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	var partFunc mrstorage.SQLPartFunc

	if fn != nil {
		partFunc = fn(b.helper)
	}

	if partFunc != nil {
		return partFunc
	}

	return b.defaultPartFunc
}

func (b *SQLOrderByBuilder) createPart(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return mrsql.NewPart(startArgumentNumber, part)
}
