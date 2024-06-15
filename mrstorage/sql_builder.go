package mrstorage

import (
	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	// SQLBuilderSet - построитель выражений используемых в SET (список поле=значение через запятую).
	SQLBuilderSet interface {
		Join(fields ...SQLBuilderPartFunc) SQLBuilderPartFunc
		Field(name string, value any) SQLBuilderPartFunc
		Fields(names []string, args []any) SQLBuilderPartFunc
	}

	// SQLBuilderWhere - построитель выражений используемых в WHERE.
	SQLBuilderWhere interface { //nolint:interfacebloat
		JoinAnd(conds ...SQLBuilderPartFunc) SQLBuilderPartFunc
		JoinOr(conds ...SQLBuilderPartFunc) SQLBuilderPartFunc

		Expr(expr string) SQLBuilderPartFunc
		ExprWithValue(expr string, value any) SQLBuilderPartFunc

		Equal(name string, value any) SQLBuilderPartFunc
		NotEqual(name string, value any) SQLBuilderPartFunc
		Less(name string, value any) SQLBuilderPartFunc
		LessOrEqual(name string, value any) SQLBuilderPartFunc
		Greater(name string, value any) SQLBuilderPartFunc
		GreaterOrEqual(name string, value any) SQLBuilderPartFunc

		FilterEqualString(name, value string) SQLBuilderPartFunc
		FilterEqualInt64(name string, value, empty int64) SQLBuilderPartFunc
		FilterEqualUUID(name string, value uuid.UUID) SQLBuilderPartFunc
		FilterEqualBool(name string, value *bool) SQLBuilderPartFunc
		FilterLike(name, value string) SQLBuilderPartFunc
		FilterLikeFields(names []string, value string) SQLBuilderPartFunc
		FilterRangeInt64(name string, value mrtype.RangeInt64, empty int64) SQLBuilderPartFunc
		// FilterAnyOf - 'values' support only slices else the func returns nil
		FilterAnyOf(name string, values any) SQLBuilderPartFunc
	}

	// SQLBuilderOrderBy - построитель выражений используемых в ORDER BY.
	SQLBuilderOrderBy interface {
		Join(fields ...SQLBuilderPartFunc) SQLBuilderPartFunc
		Field(name string, direction mrenum.SortDirection) SQLBuilderPartFunc
	}

	// SQLBuilderLimit - построитель выражений используемых в LIMIT.
	SQLBuilderLimit interface {
		OffsetLimit(index, size uint64) SQLBuilderPartFunc
	}

	// SQLBuilderCondition - помощник для построения условий объединяющий SQLBuilderWhere выражения.
	SQLBuilderCondition interface {
		Where(f func(w SQLBuilderWhere) SQLBuilderPartFunc) SQLBuilderPart
	}

	// SQLBuilderSelect - помощник для построения SELECT запросов.
	SQLBuilderSelect interface {
		SQLBuilderCondition
		OrderBy(f func(o SQLBuilderOrderBy) SQLBuilderPartFunc) SQLBuilderPart
		Limit(f func(p SQLBuilderLimit) SQLBuilderPartFunc) SQLBuilderPart
	}

	// SQLBuilderUpdate - помощник для построения UPDATE запросов.
	SQLBuilderUpdate interface {
		Set(f func(s SQLBuilderSet) SQLBuilderPartFunc) SQLBuilderPart
		SetFromEntity(entity any) (SQLBuilderPart, error)
		SetFromEntityWith(entity any, extFields func(s SQLBuilderSet) SQLBuilderPartFunc) (SQLBuilderPart, error)
		SQLBuilderCondition
	}

	// SQLSelectParams - параметры используемы при построении SELECT запросов.
	SQLSelectParams struct {
		Where   SQLBuilderPart
		OrderBy SQLBuilderPart
		Limit   SQLBuilderPart
	}
)
