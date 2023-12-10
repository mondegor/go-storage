package mrstorage

import "fmt"

type (
	SqlBuilderPart interface {
		WithPrefix(value string) SqlBuilderPart
		Param(number int) SqlBuilderPart
		Empty() bool
		ToSql() (string, []any)
		fmt.Stringer
	}

	SqlBuilderPartFunc func(paramNumber int) (string, []any)
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
