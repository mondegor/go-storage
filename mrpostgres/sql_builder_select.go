package mrpostgres

import (
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	SqlBuilderSelect struct {
		where   *SqlBuilderWhere
		orderBy *SqlBuilderOrderBy
		pager   *SqlBuilderPager
	}
)

func NewSqlBuilderSelect(
	where *SqlBuilderWhere,
	orderBy *SqlBuilderOrderBy,
	pager *SqlBuilderPager,
) *SqlBuilderSelect {
	return &SqlBuilderSelect{
		where:   where,
		orderBy: orderBy,
		pager:   pager,
	}
}

func NewSqlBuilderSelectCondition(
	where *SqlBuilderWhere,
) *SqlBuilderSelect {
	return &SqlBuilderSelect{
		where: where,
	}
}

func (b *SqlBuilderSelect) Where(f func(w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
	var partFunc mrstorage.SqlBuilderPartFunc

	if b.where != nil {
		partFunc = f(b.where)
	}

	return mrsql.NewBuilderPart(partFunc)
}

func (b *SqlBuilderSelect) OrderBy(f func(o mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
	var partFunc mrstorage.SqlBuilderPartFunc

	if b.orderBy != nil {
		partFunc = f(b.orderBy)

		if partFunc == nil {
			partFunc = b.orderBy.DefaultField()
		}
	}

	return mrsql.NewBuilderPart(partFunc)
}

func (b *SqlBuilderSelect) Pager(f func(p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
	var partFunc mrstorage.SqlBuilderPartFunc

	if b.pager != nil {
		partFunc = f(b.pager)
	}

	return mrsql.NewBuilderPart(partFunc)
}
