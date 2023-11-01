package mrstorage

import (
    "github.com/mondegor/go-storage/mrentity"
)

type (
    SqlBuilderSet interface {
        DbName(name string) string
        Join(fields ...SqlBuilderPartFunc) SqlBuilderPartFunc
        Field(dbName string, value any) SqlBuilderPartFunc
        Fields(dbNames []string, args []any) SqlBuilderPartFunc
    }

    SqlBuilderWhere interface {
        JoinAnd(conds ...SqlBuilderPartFunc) SqlBuilderPartFunc
        JoinOr(conds ...SqlBuilderPartFunc) SqlBuilderPartFunc
        Expr(expr string) SqlBuilderPartFunc
        ExprWithValue(expr string, value any) SqlBuilderPartFunc
        Equal(dbName string, value any) SqlBuilderPartFunc
        NotEqual(dbName string, value any) SqlBuilderPartFunc
        FilterLike(dbName string, value string) SqlBuilderPartFunc
        FilterLikeFields(dbNames []string, value string) SqlBuilderPartFunc
        FilterEqualInt64(dbName string, value int64, empty int64) SqlBuilderPartFunc
        FilterRangeInt64(dbName string, value mrentity.RangeInt64, empty int64) SqlBuilderPartFunc
        FilterAnyOf(dbName string, values any) SqlBuilderPartFunc
    }

    SqlBuilderOrderBy interface {
        DbName(name string) string
        Join(fields ...SqlBuilderPartFunc) SqlBuilderPartFunc
        Field(dbName string, direction mrentity.SortDirection) SqlBuilderPartFunc
    }

    SqlBuilderPager interface {
        OffsetLimit(index uint64, size uint64) SqlBuilderPartFunc
    }

    SqlBuilderSelect interface {
        Where(f func (w SqlBuilderWhere) SqlBuilderPartFunc) SqlBuilderPart
        OrderBy(f func (o SqlBuilderOrderBy) SqlBuilderPartFunc) SqlBuilderPart
        Pager(f func (p SqlBuilderPager) SqlBuilderPartFunc) SqlBuilderPart
    }

    SqlSelectParams struct {
        Where   SqlBuilderPart
        OrderBy SqlBuilderPart
        Pager   SqlBuilderPart
    }

    SqlBuilderUpdate interface {
        Set(f func(s SqlBuilderSet) SqlBuilderPartFunc) SqlBuilderPart
        SetFromEntity(entity any) (SqlBuilderPart, error)
    }
)
