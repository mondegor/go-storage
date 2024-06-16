package mrpostgres

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderSet - comment struct.
	SQLBuilderSet struct{}
)

// NewSQLBuilderSet - создаёт объект SQLBuilderSet.
func NewSQLBuilderSet() *SQLBuilderSet {
	return &SQLBuilderSet{}
}

// Join - comment method.
func (b *SQLBuilderSet) Join(fields ...mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	fields = mrstorage.SQLBuilderPartFuncRemoveNil(fields)

	if len(fields) == 0 {
		return nil
	}

	return func(_ int) (string, []any) {
		var prepared []string

		for i := range fields {
			item, _ := fields[i](0)
			prepared = append(prepared, item)
		}

		return strings.Join(prepared, ", "), nil
	}
}

// Field - comment method.
func (b *SQLBuilderSet) Field(name string, value any) mrstorage.SQLBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value}
	}
}

// Fields - comment method.
func (b *SQLBuilderSet) Fields(names []string, args []any) mrstorage.SQLBuilderPartFunc {
	if len(names) == 0 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		set := make([]string, len(names))

		for i := range names {
			set[i] = names[i] + " = $" + strconv.Itoa(paramNumber)
			paramNumber++
		}

		return strings.Join(set, ", "), args
	}
}
