package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	SQLBuilderSelect struct {
		where   *SQLBuilderWhere
		orderBy *SQLBuilderOrderBy
		pager   *SQLBuilderPager
	}
)

func NewSQLBuilderSelect(
	where *SQLBuilderWhere,
	orderBy *SQLBuilderOrderBy,
	pager *SQLBuilderPager,
) *SQLBuilderSelect {
	return &SQLBuilderSelect{
		where:   where,
		orderBy: orderBy,
		pager:   pager,
	}
}

func NewSQLBuilderSelectCondition(
	where *SQLBuilderWhere,
) *SQLBuilderSelect {
	return &SQLBuilderSelect{
		where: where,
	}
}

func (b *SQLBuilderSelect) Where(f func(w mrstorage.SQLBuilderWhere) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	}

	return mrsql.NewBuilderPart(partFunc)
}

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

func (b *SQLBuilderSelect) Pager(f func(p mrstorage.SQLBuilderPager) mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPart {
	var partFunc mrstorage.SQLBuilderPartFunc

	if b.pager != nil {
		partFunc = f(b.pager)
	}

	return mrsql.NewBuilderPart(partFunc)
}
