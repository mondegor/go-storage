package part

import (
	"github.com/mondegor/go-webcore/mrcore"

	"github.com/mondegor/go-storage/mrpostgres/builder/helper"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

const (
	startArgumentNumber = 1
)

type (
	// SQLSetBuilder - объект для создания части SQL используемой в UPDATE SET (field = $1, ...).
	// Поддерживает создание части SQL на основе тегов структуры, что помогает динамически формировать список нужных полей,
	// а также обновлять поля, которые были явно указаны.
	SQLSetBuilder struct {
		helper *helper.SQLSet
		meta   *mrsql.EntityMetaUpdate
	}
)

// NewSQLSetBuilder - создаёт объект SQLSetBuilder.
func NewSQLSetBuilder(meta *mrsql.EntityMetaUpdate) *SQLSetBuilder {
	return &SQLSetBuilder{
		helper: helper.NewSQLSet(),
		meta:   meta,
	}
}

// Build - создаёт часть SQL, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLSetBuilder) Build(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return b.createPart(part)
}

// BuildComma - создаёт часть SQL объединяющую независимые части через запятую, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLSetBuilder) BuildComma(parts ...mrstorage.SQLPartFunc) mrstorage.SQLPart {
	parts = mrsql.SQLPartFuncRemoveNil(parts)

	if len(parts) == 0 {
		return b.createPart(nil)
	}

	if len(parts) == 1 {
		return b.createPart(parts[0])
	}

	return b.createPart(b.helper.JoinComma(parts...))
}

// BuildEntity - создаёт часть SQL выбирая значения из указанной структуры, информация о которой указывается в конструкторе.
// Возвращаемое значение предназначено быть частью конкретного SQL выражения.
func (b *SQLSetBuilder) BuildEntity(entity any, parts ...mrstorage.SQLPartFunc) (mrstorage.SQLPart, error) {
	if b.meta == nil {
		return nil, mrcore.ErrInternalNilPointer.New()
	}

	dbNames, args, err := b.meta.FieldsForUpdate(entity)
	if err != nil {
		return nil, err
	}

	parts = mrsql.SQLPartFuncRemoveNil(parts)

	if len(parts) == 0 {
		return b.Build(b.helper.Fields(dbNames, args)), nil
	}

	parts = append(parts, nil) // здесь не меньше двух элементов

	// сдвиг ряда на один элемент вправо для вставки основного элемента в начало
	for i := len(parts) - 2; i >= 0; i-- {
		parts[i+1] = parts[i]
	}

	parts[0] = b.helper.Fields(dbNames, args)

	return b.Build(b.helper.JoinComma(parts...)), nil
}

// BuildFunc - создаёт часть SQL с использованием помощника, которая предназначена быть частью конкретного SQL выражения.
func (b *SQLSetBuilder) BuildFunc(fn func(s mrstorage.SQLSetHelper) mrstorage.SQLPartFunc) mrstorage.SQLPart {
	if fn != nil {
		return b.createPart(fn(b.helper))
	}

	return b.createPart(nil)
}

// HelpFunc - создаёт независимую часть SQL, которая может быть использована при создании других частей SQL.
func (b *SQLSetBuilder) HelpFunc(fn func(s mrstorage.SQLSetHelper) mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	if fn != nil {
		return fn(b.helper)
	}

	return nil
}

func (b *SQLSetBuilder) createPart(part mrstorage.SQLPartFunc) mrstorage.SQLPart {
	return mrsql.NewPart(startArgumentNumber, part)
}
