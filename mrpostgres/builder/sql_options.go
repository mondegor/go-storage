package builder

import (
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrsql"
)

type (
	// Option - настройка объекта SQL.
	Option func(o *options)

	options struct {
		sql *SQL
	}
)

// WithSQLSetMetaEntity - устанавливает для SQL метаинформацию загруженную из тегов структуры.
func WithSQLSetMetaEntity(value *mrsql.EntityMetaUpdate) Option {
	return func(o *options) {
		o.sql.set = part.NewSQLSetBuilder(value)
	}
}

// WithSQLOrderByDefaultSort - устанавливает опцию сортировка по умолчанию.
func WithSQLOrderByDefaultSort(value mrtype.SortParams) Option {
	return func(o *options) {
		o.sql.orderBy = part.NewSQLOrderByBuilder(value)
	}
}

// WithSQLLimitMaxSize - устанавливает для SQL опцию максимального кол-во строк,
// которое может быть выбрано за одни запрос.
func WithSQLLimitMaxSize(value int) Option {
	return func(o *options) {
		o.sql.limit = part.NewSQLLimitBuilder(value)
	}
}
