package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	BuilderPart struct {
		paramNumber int
		prefix      string
		bodyFunc    mrstorage.SqlBuilderPartFunc
	}
)

func NewBuilderPart(body mrstorage.SqlBuilderPartFunc) *BuilderPart {
	return &BuilderPart{
		paramNumber: 1,
		bodyFunc:    body,
	}
}

func (b *BuilderPart) Param(number int) mrstorage.SqlBuilderPart {
	if b.paramNumber == number {
		return b
	}

	c := *b
	c.paramNumber = number

	return &c
}

func (b *BuilderPart) WithPrefix(value string) mrstorage.SqlBuilderPart {
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

func (b *BuilderPart) ToSql() (string, []any) {
	return b.toSql()
}

func (b *BuilderPart) String() string {
	body, _ := b.toSql()

	return body
}

func (b *BuilderPart) toSql() (string, []any) {
	if b.bodyFunc == nil {
		return "", []any{}
	}

	body, args := b.bodyFunc(b.paramNumber)

	return b.prefix + body, args
}
