package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// BuilderPart - comment struct.
	BuilderPart struct {
		paramNumber int
		prefix      string
		bodyFunc    mrstorage.SQLBuilderPartFunc
	}
)

// NewBuilderPart - comment func.
func NewBuilderPart(body mrstorage.SQLBuilderPartFunc) *BuilderPart {
	return &BuilderPart{
		paramNumber: 1,
		bodyFunc:    body,
	}
}

// WithPart - comment method.
func (b *BuilderPart) WithPart(sep string, next mrstorage.SQLBuilderPart) mrstorage.SQLBuilderPart {
	if next == nil {
		return b
	}

	c := *b
	c.bodyFunc = func(paramNumber int) (string, []any) {
		if b.bodyFunc == nil {
			return next.WithParam(paramNumber).ToSQL()
		}

		body1, args1 := b.bodyFunc(paramNumber)
		body2, args2 := next.WithParam(paramNumber + len(args1)).ToSQL()

		return b.prefix + body1 + sep + body2, MergeArgs(args1, args2)
	}

	return &c
}

// WithParam - comment method.
func (b *BuilderPart) WithParam(number int) mrstorage.SQLBuilderPart {
	if b.paramNumber == number {
		return b
	}

	c := *b
	c.paramNumber = number

	return &c
}

// WithPrefix - comment method.
func (b *BuilderPart) WithPrefix(value string) mrstorage.SQLBuilderPart {
	if b.prefix == value {
		return b
	}

	c := *b
	c.prefix = value

	return &c
}

// Empty - проверяет, что в объекте не установлена функция для формирования SQL.
func (b *BuilderPart) Empty() bool {
	return b.bodyFunc == nil
}

// ToSQL - comment method.
func (b *BuilderPart) ToSQL() (string, []any) {
	return b.toSQL()
}

// String - comment method.
func (b *BuilderPart) String() string {
	body, _ := b.toSQL()

	return body
}

func (b *BuilderPart) toSQL() (string, []any) {
	if b.bodyFunc == nil {
		return "", nil
	}

	body, args := b.bodyFunc(b.paramNumber)

	return b.prefix + body, args
}
