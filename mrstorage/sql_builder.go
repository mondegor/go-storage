package mrstorage

import (
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	SqlBuilderSet interface {
		Join(fields ...SqlBuilderPartFunc) SqlBuilderPartFunc
		Field(name string, value any) SqlBuilderPartFunc
		Fields(names []string, args []any) SqlBuilderPartFunc
	}

	SqlBuilderWhere interface {
		JoinAnd(conds ...SqlBuilderPartFunc) SqlBuilderPartFunc
		JoinOr(conds ...SqlBuilderPartFunc) SqlBuilderPartFunc
		Expr(expr string) SqlBuilderPartFunc
		ExprWithValue(expr string, value any) SqlBuilderPartFunc
		Equal(name string, value any) SqlBuilderPartFunc
		NotEqual(name string, value any) SqlBuilderPartFunc
		FilterEqualString(name, value string) SqlBuilderPartFunc
		FilterEqualInt64(name string, value, empty int64) SqlBuilderPartFunc
		FilterEqualBool(name string, value mrtype.NullableBool) SqlBuilderPartFunc
		FilterLike(name, value string) SqlBuilderPartFunc
		FilterLikeFields(names []string, value string) SqlBuilderPartFunc
		FilterRangeInt64(name string, value mrtype.RangeInt64, empty int64) SqlBuilderPartFunc
		FilterAnyOf(name string, values any) SqlBuilderPartFunc
	}

	SqlBuilderOrderBy interface {
		WrapWithDefault(field SqlBuilderPartFunc) SqlBuilderPartFunc
		Join(fields ...SqlBuilderPartFunc) SqlBuilderPartFunc
		Field(name string, direction mrenum.SortDirection) SqlBuilderPartFunc
	}

	SqlBuilderPager interface {
		OffsetLimit(index, size uint64) SqlBuilderPartFunc
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
		SetFromEntityWith(entity any, extFields func(s SqlBuilderSet) SqlBuilderPartFunc) (SqlBuilderPart, error)
	}
)
