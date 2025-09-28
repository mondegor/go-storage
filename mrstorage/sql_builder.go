package mrstorage

import (
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/enums"
)

type (
	// SQLBuilder - строитель условий используемых в конструкции SET.
	SQLBuilder interface {
		Set() SQLSetBuilder
		Condition() SQLConditionBuilder
		OrderBy() SQLOrderByBuilder
		Limit() SQLLimitBuilder
	}

	// SQLPart - параметризованная часть SQL запроса.
	SQLPart interface {
		WithPrefix(sql string) SQLPart
		WithStartArg(number int) SQLPart
		Empty() bool
		String() string
		ToSQL() (sql string, args []any)
	}

	// SQLSetBuilder - строитель условий используемых в конструкции SET.
	SQLSetBuilder interface {
		Build(part SQLPartFunc) SQLPart
		BuildComma(parts ...SQLPartFunc) SQLPart
		BuildEntity(entity any, parts ...SQLPartFunc) (SQLPart, error)
		BuildFunc(fn func(s SQLSetHelper) SQLPartFunc) SQLPart
		HelpFunc(fn func(s SQLSetHelper) SQLPartFunc) SQLPartFunc
	}

	// SQLSetHelper - помощник для построения выражений используемых в конструкции SET.
	SQLSetHelper interface {
		JoinComma(fields ...SQLPartFunc) SQLPartFunc
		Field(name string, value any) SQLPartFunc
		Fields(names []string, args []any) SQLPartFunc
	}

	// SQLConditionBuilder - строитель условий используемых в WHERE, JOIN конструкциях.
	SQLConditionBuilder interface {
		Build(part SQLPartFunc) SQLPart
		BuildAnd(parts ...SQLPartFunc) SQLPart
		BuildFunc(fn func(c SQLConditionHelper) SQLPartFunc) SQLPart
		HelpFunc(fn func(c SQLConditionHelper) SQLPartFunc) SQLPartFunc
	}

	// SQLConditionHelper - помощник для построения выражений используемых в WHERE, JOIN конструкциях.
	SQLConditionHelper interface {
		JoinAnd(parts ...SQLPartFunc) SQLPartFunc
		JoinOr(parts ...SQLPartFunc) SQLPartFunc

		Expr(expr string) SQLPartFunc
		ExprWithValue(expr string, value any) SQLPartFunc

		Equal(field string, value any) SQLPartFunc
		NotEqual(field string, value any) SQLPartFunc
		Less(field string, value any) SQLPartFunc
		LessOrEqual(field string, value any) SQLPartFunc
		Greater(field string, value any) SQLPartFunc
		GreaterOrEqual(field string, value any) SQLPartFunc

		FilterEqual(field string, value any) SQLPartFunc
		FilterEqualString(field, value string) SQLPartFunc
		FilterEqualInt64(field string, value, empty int64) SQLPartFunc
		FilterEqualBool(field string, value *bool) SQLPartFunc
		FilterLike(field, value string) SQLPartFunc
		FilterLikeFields(fields []string, value string) SQLPartFunc
		FilterRangeInt64(field string, value mrtype.RangeInt64, empty int64) SQLPartFunc
		FilterRangeFloat64(field string, value mrtype.RangeFloat64, empty, qualityThreshold float64) SQLPartFunc
		// FilterAnyOf - used ANY(), 'values' support only slices else the func returns nil
		FilterAnyOf(field string, values any) SQLPartFunc
		// FilterInOf - used IN(), 'values' support only slices else the func returns nil
		FilterInOf(field string, values any) SQLPartFunc
	}

	// SQLOrderByBuilder - строитель условий используемых в конструкции ORDER BY.
	SQLOrderByBuilder interface {
		Build(part SQLPartFunc) SQLPart
		BuildComma(parts ...SQLPartFunc) SQLPart
		BuildFunc(fn func(o SQLOrderByHelper) SQLPartFunc) SQLPart
		HelpFunc(fn func(o SQLOrderByHelper) SQLPartFunc) SQLPartFunc
	}

	// SQLOrderByHelper - помощник для построения выражений используемых в конструкции ORDER BY.
	SQLOrderByHelper interface {
		JoinComma(fields ...SQLPartFunc) SQLPartFunc
		Field(name string, direction enums.SortDirection) SQLPartFunc
	}

	// SQLLimitBuilder - строитель условий используемых в конструкции LIMIT.
	SQLLimitBuilder interface {
		Build(index, size uint64) SQLPart
	}

	// SQLPartFunc - динамическая часть SQL запроса, вычисляемая тогда,
	// когда понятен номер первого параметра используемого в этой функции.
	SQLPartFunc func(argumentNumber int) (sql string, args []any)
)
