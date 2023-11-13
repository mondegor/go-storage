package mrsql

import (
	"fmt"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	BuilderPart struct {
		paramNumber int
		prefix string
		bodyFunc mrstorage.SqlBuilderPartFunc
	}
)

func NewBuilderPart(body mrstorage.SqlBuilderPartFunc) *BuilderPart {
	return &BuilderPart{
		paramNumber: 1,
		bodyFunc: body,
	}
}

func (b *BuilderPart) Empty() bool {
	return b.bodyFunc == nil
}

func (b *BuilderPart) WithPrefix(value string) mrstorage.SqlBuilderPart {
	return &BuilderPart{
		paramNumber: b.paramNumber,
		prefix: value,
		bodyFunc: b.bodyFunc,
	}
}

func (b *BuilderPart) Param(number int) mrstorage.SqlBuilderPart {
	return &BuilderPart{
		paramNumber: number,
		prefix: b.prefix,
		bodyFunc: b.bodyFunc,
	}
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

	return fmt.Sprintf("%s%s", b.prefix, body), args
}
