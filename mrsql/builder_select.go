package mrsql

import "github.com/mondegor/go-storage/mrstorage"

type (
    BuilderSelect struct {
        where mrstorage.SqlBuilderWhere
        orderBy mrstorage.SqlBuilderOrderBy
        pager mrstorage.SqlBuilderPager
    }
)

func NewBuilderSelect(
    where mrstorage.SqlBuilderWhere,
    orderBy mrstorage.SqlBuilderOrderBy,
    pager mrstorage.SqlBuilderPager,
) *BuilderSelect {
    return &BuilderSelect{
        where: where,
        orderBy: orderBy,
        pager: pager,
    }
}

func (b *BuilderSelect) Where(f func (w mrstorage.SqlBuilderWhere) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
    return NewBuilderPart(f(b.where))
}

func (b *BuilderSelect) OrderBy(f func (o mrstorage.SqlBuilderOrderBy) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
    return NewBuilderPart(
        b.orderBy.WrapWithDefault(
            f(b.orderBy),
        ),
    )
}

func (b *BuilderSelect) Pager(f func (p mrstorage.SqlBuilderPager) mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPart {
    return NewBuilderPart(f(b.pager))
}
