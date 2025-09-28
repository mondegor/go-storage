package builder

import (
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrsql"
)

type (
	// Option - настройка объекта SQL.
	Option func(b *SQL)
)

// WithSQLSetMetaEntity - устанавливает для SQL метаинформацию загруженную из тегов структуры.
func WithSQLSetMetaEntity(value *mrsql.EntityMetaUpdate) Option {
	return func(b *SQL) {
		b.set = part.NewSQLSetBuilder(value)
	}
}

// WithSQLOrderByDefaultSort - устанавливает опцию сортировка по умолчанию.
func WithSQLOrderByDefaultSort(value mrtype.SortParams) Option {
	return func(b *SQL) {
		b.orderBy = part.NewSQLOrderByBuilder(value)
	}
}

// WithSQLLimitMaxSize - устанавливает для SQL опцию максимального кол-во строк,
// которое может быть выбрано за одни запрос.
func WithSQLLimitMaxSize(value uint64) Option {
	return func(b *SQL) {
		b.limit = part.NewSQLLimitBuilder(value)
	}
}
