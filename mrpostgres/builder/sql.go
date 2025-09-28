package builder

import (
	"github.com/mondegor/go-sysmess/mrtype"

	"github.com/mondegor/go-storage/mrpostgres/builder/part"
	"github.com/mondegor/go-storage/mrstorage"
)

const (
	defaultLimitMaxSize = 5
)

type (
	// SQL - агрегирующий вспомогательный объект для создания различных частей SQL выражения.
	// Данный объект не является универсальным, он создаётся под конкретный репозиторий
	// для решения конкретных задач этого репозитория.
	SQL struct {
		set       *part.SQLSetBuilder
		condition *part.SQLConditionBuilder
		orderBy   *part.SQLOrderByBuilder
		limit     *part.SQLLimitBuilder
	}
)

// NewSQL - создаёт объект SQL.
func NewSQL(opts ...Option) *SQL {
	b := &SQL{
		set:       part.NewSQLSetBuilder(nil), // WARNING: по умолчанию EntityMeta не указано
		condition: part.NewSQLConditionBuilder(),
		orderBy:   part.NewSQLOrderByBuilder(mrtype.SortParams{}), // сортировка по умолчанию отсутствует
		limit:     part.NewSQLLimitBuilder(defaultLimitMaxSize),
	}

	for _, opt := range opts {
		opt(b)
	}

	return b
}

// Set - возвращает объект для создания части SQL используемой в UPDATE SET (field = $1, ...).
func (b *SQL) Set() mrstorage.SQLSetBuilder {
	return b.set
}

// Condition - возвращает объект для создания части SQL используемой в WHERE, JOIN (field = $1 AND ...).
func (b *SQL) Condition() mrstorage.SQLConditionBuilder {
	return b.condition
}

// OrderBy - возвращает объект для создания части SQL используемой в ORDER BY (field ASC, ...).
func (b *SQL) OrderBy() mrstorage.SQLOrderByBuilder {
	return b.orderBy
}

// Limit - возвращает объект для создания части SQL используемой в OFFSET, LIMIT.
func (b *SQL) Limit() mrstorage.SQLLimitBuilder {
	return b.limit
}
