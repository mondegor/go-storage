package part

import (
	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLConditionBuilder - объект для создания части SQL используемой в WHERE, JOIN (field = $1 AND ...).
	SQLConditionBuilder struct {
		helper *helper.SQLCondition
	}
)

// NewSQLConditionBuilder - создаёт объект SQLConditionBuilder.
func NewSQLConditionBuilder() *SQLConditionBuilder {
	return &SQLConditionBuilder{
		helper: helper.NewSQLCondition(),
	}
}

// Build - создаёт часть SQL, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLConditionBuilder) Build(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return b.createPart(part)
}

// BuildAnd - создаёт часть SQL объединяющую независимые части через оператор AND, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLConditionBuilder) BuildAnd(parts ...mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return b.createPart(b.helper.JoinAnd(parts...))
}

// BuildFunc - создаёт часть SQL с использованием помощника, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLConditionBuilder) BuildFunc(fn func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc) mrstorage.SQLPart {
	if fn != nil {
		return b.createPart(fn(b.helper))
	}

	return b.createPart(nil)
}

// HelpFunc - создаёт независимую часть SQL, которая может быть использована при создании других частей SQL.
func (b *SQLConditionBuilder) HelpFunc(fn func(c mrstorage.SQLConditionHelper) mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	if fn != nil {
		return fn(b.helper)
	}

	return nil
}

func (b *SQLConditionBuilder) createPart(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return mrsql.NewPart(startArgumentNumber, part)
}
