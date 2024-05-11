package mrstorage

import "fmt"

type (
	SQLBuilderPart interface {
		WithPrefix(value string) SQLBuilderPart
		Param(number int) SQLBuilderPart
		Empty() bool
		ToSQL() (string, []any)
		fmt.Stringer
	}

	SQLBuilderPartFunc func(paramNumber int) (string, []any)
)

func SQLBuilderPartFuncRemoveNil(parts []SQLBuilderPartFunc) []SQLBuilderPartFunc {
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
