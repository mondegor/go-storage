package mrpostgres

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	SqlBuilderSet struct{}
)

func NewSqlBuilderSet() *SqlBuilderSet {
	return &SqlBuilderSet{}
}

func (b *SqlBuilderSet) Join(fields ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	fields = mrstorage.SqlBuilderPartFuncRemoveNil(fields)

	if len(fields) == 0 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		var prepared []string

		for i := range fields {
			item, _ := fields[i](0)
			prepared = append(prepared, item)
		}

		return strings.Join(prepared, ", "), []any{}
	}
}

func (b *SqlBuilderSet) Field(name string, value any) mrstorage.SqlBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value}
	}
}

func (b *SqlBuilderSet) Fields(names []string, args []any) mrstorage.SqlBuilderPartFunc {
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
