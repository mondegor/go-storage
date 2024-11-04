package helper

import (
	"strconv"
	"strings"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLSet - объект для создания независимой части SQL используемой в присвоении значений полей (SET).
	SQLSet struct{}
)

// NewSQLSet - создаёт объект SQLSet.
func NewSQLSet() *SQLSet {
	return &SQLSet{}
}

// JoinComma - возвращает указанные SQL поля соединённые через запятую.
func (b *SQLSet) JoinComma(fields ...mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
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

// Field - возвращает SQL поле с присвоенным ему значением.
func (b *SQLSet) Field(name string, value any) mrstorage.SQLPartFunc {
	return func(argumentNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(argumentNumber), []any{value}
	}
}

// Fields - возвращает SQL поля с присвоенными им значениями.
func (b *SQLSet) Fields(names []string, args []any) mrstorage.SQLPartFunc {
	if len(names) == 0 {
		return nil
	}

	return func(argumentNumber int) (string, []any) {
		set := make([]string, len(names))

		for i := range names {
			set[i] = names[i] + " = $" + strconv.Itoa(argumentNumber)
			argumentNumber++
		}

		return strings.Join(set, ", "), args
	}
}
