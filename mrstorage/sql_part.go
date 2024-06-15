package mrstorage

type (
	// SQLBuilderPart - построитель запросов.
	SQLBuilderPart interface {
		WithPart(sep string, next SQLBuilderPart) SQLBuilderPart
		WithPrefix(value string) SQLBuilderPart
		WithParam(number int) SQLBuilderPart
		Empty() bool
		ToSQL() (string, []any)
		String() string
	}

	// SQLBuilderPartFunc - часть запроса зависящая от параметров.
	SQLBuilderPartFunc func(paramNumber int) (string, []any)
)

// SQLBuilderPartFuncRemoveNil - удаляет все nil элементы из указанного массива.
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
