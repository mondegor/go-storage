package mrstorage

import "fmt"

type (
	SqlBuilderPart interface {
		Empty() bool
		Param(number int) SqlBuilderPart
		WithPrefix(value string) SqlBuilderPart
		ToSql() (string, []any)
		fmt.Stringer
	}

	SqlBuilderPartFunc func (paramNumber int) (string, []any)
)

func SqlBuilderPartFuncRemoveNil(parts []SqlBuilderPartFunc) []SqlBuilderPartFunc {
	needOffset := false
	length := 0

	for i := range parts {
		if parts[i] == nil {
			needOffset = true
			continue
		}

		if i > length {
			parts[length] = parts[i]
		}

		length++
	}

	if needOffset {
		parts = parts[:length]
	}

	return parts
}
