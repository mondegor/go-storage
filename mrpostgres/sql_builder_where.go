package mrpostgres

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/lib/pq"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

// go get -u github.com/lib/pq

type (
	SqlBuilderWhere struct {
	}
)

func NewSqlBuilderWhere() *SqlBuilderWhere {
	return &SqlBuilderWhere{}
}

func (b *SqlBuilderWhere) JoinAnd(conds ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	return b.join(" AND ", conds)
}

func (b *SqlBuilderWhere) JoinOr(conds ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	return b.join(" OR ", conds)
}

func (b *SqlBuilderWhere) Expr(expr string) mrstorage.SqlBuilderPartFunc {
	if expr == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return expr, []any{}
	}
}

// ExprWithValue - sample: "UPPER(field_name) = %s"
func (b *SqlBuilderWhere) ExprWithValue(expr string, value any) mrstorage.SqlBuilderPartFunc {
	if expr == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return fmt.Sprintf(expr, "$"+strconv.Itoa(paramNumber)), []any{value}
	}
}

func (b *SqlBuilderWhere) Equal(name string, value any) mrstorage.SqlBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value}
	}
}

func (b *SqlBuilderWhere) NotEqual(name string, value any) mrstorage.SqlBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " <> $" + strconv.Itoa(paramNumber), []any{value}
	}
}

func (b *SqlBuilderWhere) FilterEqualString(name, value string) mrstorage.SqlBuilderPartFunc {
	if value == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value}
	}
}

func (b *SqlBuilderWhere) FilterEqualInt64(name string, value, empty int64) mrstorage.SqlBuilderPartFunc {
	if value == empty {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value}
	}
}

func (b *SqlBuilderWhere) FilterEqualBool(name string, value mrtype.NullableBool) mrstorage.SqlBuilderPartFunc {
	if value.IsNull() {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " = $" + strconv.Itoa(paramNumber), []any{value.Val()}
	}
}

func (b *SqlBuilderWhere) FilterLike(name, value string) mrstorage.SqlBuilderPartFunc {
	return b.FilterLikeFields([]string{name}, value)
}

func (b *SqlBuilderWhere) FilterLikeFields(names []string, value string) mrstorage.SqlBuilderPartFunc {
	if value == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		var conds []string

		for i := range names {
			conds = append(conds, names[i]+" LIKE '%%' || $"+strconv.Itoa(paramNumber)+" || '%%'")
		}

		return "(" + strings.Join(conds, " OR ") + ")", []any{value}
	}
}

func (b *SqlBuilderWhere) FilterRangeInt64(name string, value mrtype.RangeInt64, empty int64) mrstorage.SqlBuilderPartFunc {
	if value.Min != empty {
		if value.Max != empty {
			if value.Min > value.Max {
				return nil
			}

			return func(paramNumber int) (string, []any) {
				return "(" + name + " BETWEEN $" + strconv.Itoa(paramNumber) + " AND $" + strconv.Itoa(paramNumber+1) + ")", []any{value.Min, value.Max}
			}
		} else {
			return func(paramNumber int) (string, []any) {
				return name + " >= $" + strconv.Itoa(paramNumber), []any{value.Min}
			}
		}
	} else if value.Max != empty {
		return func(paramNumber int) (string, []any) {
			return name + " <= $" + strconv.Itoa(paramNumber), []any{value.Max}
		}
	}

	return nil
}

func (b *SqlBuilderWhere) FilterAnyOf(name string, values any) mrstorage.SqlBuilderPartFunc {
	val := reflect.ValueOf(values)

	if val.Kind() != reflect.Slice || val.Len() == 0 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " = ANY($" + strconv.Itoa(paramNumber) + ")", []any{pq.Array(values)}
	}
}

func (b *SqlBuilderWhere) join(separator string, conds []mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	conds = mrstorage.SqlBuilderPartFuncRemoveNil(conds)

	if len(conds) == 0 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		var prepared []string
		var args []any

		for i := range conds {
			item, itemArgs := conds[i](paramNumber + len(args))
			prepared = append(prepared, item)
			args = mrsql.MergeArgs(args, itemArgs)
		}

		return "(" + strings.Join(prepared, separator) + ")", args
	}
}
