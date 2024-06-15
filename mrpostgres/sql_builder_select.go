package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderSelect - comment struct.
	SQLBuilderSelect struct {
		where   *SQLBuilderWhere
		orderBy *SQLBuilderOrderBy
		limit   *SQLBuilderLimit
	}
)

// NewSQLBuilderSelect - comment func.
func NewSQLBuilderSelect(
	where *SQLBuilderWhere,
	orderBy *SQLBuilderOrderBy,
	limit *SQLBuilderLimit,
) *SQLBuilderSelect {
	return &SQLBuilderSelect{
		where:   where,
		orderBy: orderBy,
		limit:   limit,
	}
}

// Where - comment method.
func (b *SQLBuilderSelect) Where(f func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	}

	return mrsql.NewBuilderPart(partFunc)
}

// OrderBy - comment method.
func (b *SQLBuilderSelect) OrderBy(f func(o mrstorage.SQLBuilderOrderBy) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.orderBy != nil {
		partFunc = f(b.orderBy)

		if partFunc == nil {
			partFunc = b.orderBy.DefaultField()
		}
	}

	return mrsql.NewBuilderPart(partFunc)
}

// Limit - comment method.
func (b *SQLBuilderSelect) Limit(f func(p mrstorage.SQLBuilderLimit) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.limit != nil {
		partFunc = f(b.limit)
	}

	return mrsql.NewBuilderPart(partFunc)
}
