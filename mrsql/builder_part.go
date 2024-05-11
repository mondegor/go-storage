package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	BuilderPart struct {
		paramNumber int
		prefix      string
		bodyFunc    mrstorage.SQLBuilderPartFunc
	}
)

func NewBuilderPart(body mrstorage.SQLBuilderPartFunc) *BuilderPart {
	return &BuilderPart{
		paramNumber: 1,
		bodyFunc:    body,
	}
}

func (b *BuilderPart) Param(number int) mrstorage.SQLBuilderPart {
	if b.paramNumber == number {
		return b
	}

	c := *b
	c.paramNumber = number

	return &c
}

func (b *BuilderPart) WithPrefix(value string) mrstorage.SQLBuilderPart {
	if b.prefix == value {
		return b
	}

	c := *b
	c.prefix = value

	return &c
}

func (b *BuilderPart) Empty() bool {
	return b.bodyFunc == nil
}

func (b *BuilderPart) ToSQL() (string, []any) {
	return b.toSQL()
}

func (b *BuilderPart) String() string {
	body, _ := b.toSQL()

	return body
}

func (b *BuilderPart) toSQL() (string, []any) {
	if b.bodyFunc == nil {
		return "", []any{}
	}

	body, args := b.bodyFunc(b.paramNumber)

	return b.prefix + body, args
}
