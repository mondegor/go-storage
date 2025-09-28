package helper

import (
	"strings"

	"github.com/mondegor/go-sysmess/mrtype/enums"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLOrderBy - объект для создания независимой части SQL используемой в порядке следования (ORDER BY).
	SQLOrderBy struct{}
)

// NewSQLOrderBy - создаёт объект SQLOrderBy.
func NewSQLOrderBy() *SQLOrderBy {
	return &SQLOrderBy{}
}

// JoinComma - возвращает указанные SQL поля соединённые через запятую.
func (b *SQLOrderBy) JoinComma(fields ...mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	fields = mrsql.SQLPartFuncRemoveNil(fields)

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

// Field - возвращает SQL поле с указанием направления сортировки.
func (b *SQLOrderBy) Field(name string, direction enums.SortDirection) mrstorage.SQLPartFunc {
	if name == "" {
		return nil
	}

	return func(_ int) (string, []any) {
		return name + " " + direction.String(), nil
	}
}
