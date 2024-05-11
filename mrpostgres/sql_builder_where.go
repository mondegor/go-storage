package mrpostgres

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	SqlBuilderWhere struct{}
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
	return b.compare(name, value, "=")
}

func (b *SqlBuilderWhere) NotEqual(name string, value any) mrstorage.SqlBuilderPartFunc {
	return b.compare(name, value, "<>")
}

func (b *SqlBuilderWhere) Less(name string, value any) mrstorage.SqlBuilderPartFunc {
	return b.compare(name, value, "<")
}

func (b *SqlBuilderWhere) LessOrEqual(name string, value any) mrstorage.SqlBuilderPartFunc {
	return b.compare(name, value, "<=")
}

func (b *SqlBuilderWhere) Greater(name string, value any) mrstorage.SqlBuilderPartFunc {
	return b.compare(name, value, ">")
}

func (b *SqlBuilderWhere) GreaterOrEqual(name string, value any) mrstorage.SqlBuilderPartFunc {
	return b.compare(name, value, ">=")
}

func (b *SqlBuilderWhere) FilterEqualString(name, value string) mrstorage.SqlBuilderPartFunc {
	if value == "" {
		return nil
	}

	return b.compare(name, value, "=")
}

func (b *SqlBuilderWhere) FilterEqualInt64(name string, value, empty int64) mrstorage.SqlBuilderPartFunc {
	if value == empty {
		return nil
	}

	return b.compare(name, value, "=")
}

func (b *SqlBuilderWhere) FilterEqualUUID(name string, value uuid.UUID) mrstorage.SqlBuilderPartFunc {
	if value == uuid.Nil {
		return nil
	}

	return b.compare(name, value, "=")
}

func (b *SqlBuilderWhere) FilterEqualBool(name string, value *bool) mrstorage.SqlBuilderPartFunc {
	if value == nil {
		return nil
	}

	return b.compare(name, *value, "=")
}

func (b *SqlBuilderWhere) FilterLike(name, value string) mrstorage.SqlBuilderPartFunc {
	return b.FilterLikeFields([]string{name}, value)
}

func (b *SqlBuilderWhere) FilterLikeFields(names []string, value string) mrstorage.SqlBuilderPartFunc {
	if value == "" {
		return nil
	}

	// sample: (field_name LIKE '%%' || $1 || '%%' OR ...)
	return func(paramNumber int) (string, []any) {
		var buf strings.Builder

		buf.Grow(30 * len(names))
		buf.WriteByte('(')

		for i := range names {
			if i > 0 {
				buf.WriteString(" OR ")
			}

			buf.WriteString(names[i])
			buf.WriteString(" LIKE '%%' || $")
			buf.WriteString(strconv.Itoa(paramNumber))
			buf.WriteString(" || '%%'")
		}

		buf.WriteByte(')')

		return buf.String(), []any{value}
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
			return b.compare(name, value.Min, ">=")
		}
	} else if value.Max != empty {
		return b.compare(name, value.Max, "<=")
	}

	return nil
}

// FilterAnyOf - 'values' support only slices else the func returns nil
func (b *SqlBuilderWhere) FilterAnyOf(name string, values any) mrstorage.SqlBuilderPartFunc {
	s := reflect.ValueOf(values)

	if s.Kind() != reflect.Slice || s.Len() < 1 {
		return nil
	}

	args := make([]any, s.Len())

	for i := range args {
		args[i] = s.Index(i).Interface()
	}

	// sample: field_name IN($1, $2, ...)
	return func(paramNumber int) (string, []any) {
		var buf strings.Builder

		buf.Grow(len(name) + 4 + 3*len(args)) // len(name) + " IN(" + "$N," * len(args) - 1
		buf.WriteString(name)
		buf.WriteString(" IN(")

		for i := range args {
			if i > 0 {
				buf.WriteByte(',')
			}

			buf.WriteByte('$')
			buf.WriteString(strconv.Itoa(paramNumber + i))
		}

		buf.WriteByte(')')

		return buf.String(), args
	}
}

func (b *SqlBuilderWhere) join(separator string, conds []mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	conds = mrstorage.SqlBuilderPartFuncRemoveNil(conds)

	if len(conds) == 0 {
		return nil
	}

	// sample: (cond1 AND cond2 AND ...)
	return func(paramNumber int) (string, []any) {
		var buf strings.Builder
		var args []any

		buf.WriteByte('(')

		for i := range conds {
			if i > 0 {
				buf.WriteString(separator)
			}

			item, itemArgs := conds[i](paramNumber + len(args))
			buf.WriteString(item)
			args = mrsql.MergeArgs(args, itemArgs)
		}

		buf.WriteByte(')')

		return buf.String(), args
	}
}

func (b *SqlBuilderWhere) compare(name string, value any, sign string) mrstorage.SqlBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " " + sign + " $" + strconv.Itoa(paramNumber), []any{value}
	}
}
