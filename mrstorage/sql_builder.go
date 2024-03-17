package mrstorage

import (
	"github.com/google/uuid"
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
		Less(name string, value any) SqlBuilderPartFunc
		LessOrEqual(name string, value any) SqlBuilderPartFunc
		Greater(name string, value any) SqlBuilderPartFunc
		GreaterOrEqual(name string, value any) SqlBuilderPartFunc

		FilterEqualString(name, value string) SqlBuilderPartFunc
		FilterEqualInt64(name string, value, empty int64) SqlBuilderPartFunc
		FilterEqualUUID(name string, value uuid.UUID) SqlBuilderPartFunc
		FilterEqualBool(name string, value *bool) SqlBuilderPartFunc
		FilterLike(name, value string) SqlBuilderPartFunc
		FilterLikeFields(names []string, value string) SqlBuilderPartFunc
		FilterRangeInt64(name string, value mrtype.RangeInt64, empty int64) SqlBuilderPartFunc
		// FilterAnyOf - 'values' support only slices else the func returns nil
		FilterAnyOf(name string, values any) SqlBuilderPartFunc
	}

	SqlBuilderOrderBy interface {
		Join(fields ...SqlBuilderPartFunc) SqlBuilderPartFunc
		Field(name string, direction mrenum.SortDirection) SqlBuilderPartFunc
	}

	SqlBuilderPager interface {
		OffsetLimit(index, size uint64) SqlBuilderPartFunc
	}

	SqlBuilderCondition interface {
		Where(f func(w SqlBuilderWhere) SqlBuilderPartFunc) SqlBuilderPart
	}

	SqlBuilderSelect interface {
		SqlBuilderCondition
		OrderBy(f func(o SqlBuilderOrderBy) SqlBuilderPartFunc) SqlBuilderPart
		Pager(f func(p SqlBuilderPager) SqlBuilderPartFunc) SqlBuilderPart
	}

	SqlBuilderUpdate interface {
		Set(f func(s SqlBuilderSet) SqlBuilderPartFunc) SqlBuilderPart
		SetFromEntity(entity any) (SqlBuilderPart, error)
		SetFromEntityWith(entity any, extFields func(s SqlBuilderSet) SqlBuilderPartFunc) (SqlBuilderPart, error)
		SqlBuilderCondition
	}

	SqlSelectParams struct {
		Where   SqlBuilderPart
		OrderBy SqlBuilderPart
		Pager   SqlBuilderPart
	}
)
