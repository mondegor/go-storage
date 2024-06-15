package mrpostgres

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/google/uuid"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderWhere - comment struct.
	SQLBuilderWhere struct{}
)

// NewSQLBuilderWhere - comment func.
func NewSQLBuilderWhere() *SQLBuilderWhere {
	return &SQLBuilderWhere{}
}

// JoinAnd - comment method.
func (b *SQLBuilderWhere) JoinAnd(conds ...mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	return b.join(" AND ", conds)
}

// JoinOr - comment method.
func (b *SQLBuilderWhere) JoinOr(conds ...mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	return b.join(" OR ", conds)
}

// Expr - comment method.
func (b *SQLBuilderWhere) Expr(expr string) mrstorage.SQLBuilderPartFunc {
	if expr == "" {
		return nil
	}

	return func(_ int) (string, []any) {
		return expr, nil
	}
}

// ExprWithValue - sample: "UPPER(field_name) = %s".
func (b *SQLBuilderWhere) ExprWithValue(expr string, value any) mrstorage.SQLBuilderPartFunc {
	if expr == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return fmt.Sprintf(expr, "$"+strconv.Itoa(paramNumber)), []any{value}
	}
}

// Equal - comment method.
func (b *SQLBuilderWhere) Equal(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, "=")
}

// NotEqual - comment method.
func (b *SQLBuilderWhere) NotEqual(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, "<>")
}

// Less - comment method.
func (b *SQLBuilderWhere) Less(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, "<")
}

// LessOrEqual - comment method.
func (b *SQLBuilderWhere) LessOrEqual(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, "<=")
}

// Greater - comment method.
func (b *SQLBuilderWhere) Greater(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, ">")
}

// GreaterOrEqual - comment method.
func (b *SQLBuilderWhere) GreaterOrEqual(name string, value any) mrstorage.SQLBuilderPartFunc {
	return b.compare(name, value, ">=")
}

// FilterEqualString - comment method.
func (b *SQLBuilderWhere) FilterEqualString(name, value string) mrstorage.SQLBuilderPartFunc {
	if value == "" {
		return nil
	}

	return b.compare(name, value, "=")
}

// FilterEqualInt64 - comment method.
func (b *SQLBuilderWhere) FilterEqualInt64(name string, value, empty int64) mrstorage.SQLBuilderPartFunc {
	if value == empty {
		return nil
	}

	return b.compare(name, value, "=")
}

// FilterEqualUUID - comment method.
func (b *SQLBuilderWhere) FilterEqualUUID(name string, value uuid.UUID) mrstorage.SQLBuilderPartFunc {
	if value == uuid.Nil {
		return nil
	}

	return b.compare(name, value, "=")
}

// FilterEqualBool - comment method.
func (b *SQLBuilderWhere) FilterEqualBool(name string, value *bool) mrstorage.SQLBuilderPartFunc {
	if value == nil {
		return nil
	}

	return b.compare(name, *value, "=")
}

// FilterLike - comment method.
func (b *SQLBuilderWhere) FilterLike(name, value string) mrstorage.SQLBuilderPartFunc {
	return b.FilterLikeFields([]string{name}, value)
}

// FilterLikeFields - comment method.
func (b *SQLBuilderWhere) FilterLikeFields(names []string, value string) mrstorage.SQLBuilderPartFunc {
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

// FilterRangeInt64 - comment method.
func (b *SQLBuilderWhere) FilterRangeInt64(name string, value mrtype.RangeInt64, empty int64) mrstorage.SQLBuilderPartFunc {
	if value.Min != empty {
		if value.Max != empty {
			if value.Min > value.Max {
				return nil
			}

			return func(paramNumber int) (string, []any) {
				return "(" + name + " BETWEEN $" + strconv.Itoa(paramNumber) + " AND $" + strconv.Itoa(paramNumber+1) + ")", []any{value.Min, value.Max}
			}
		}

		return b.compare(name, value.Min, ">=")
	} else if value.Max != empty {
		return b.compare(name, value.Max, "<=")
	}

	return nil
}

// FilterAnyOf - 'values' support only slices else the func returns nil.
func (b *SQLBuilderWhere) FilterAnyOf(name string, values any) mrstorage.SQLBuilderPartFunc {
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

func (b *SQLBuilderWhere) join(separator string, conds []mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	conds = mrstorage.SQLBuilderPartFuncRemoveNil(conds)

	if len(conds) == 0 {
		return nil
	}

	// sample: (cond1 AND cond2 AND ...)
	return func(paramNumber int) (string, []any) {
		var (
			buf  strings.Builder
			args []any
		)

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

func (b *SQLBuilderWhere) compare(name string, value any, sign string) mrstorage.SQLBuilderPartFunc {
	return func(paramNumber int) (string, []any) {
		return name + " " + sign + " $" + strconv.Itoa(paramNumber), []any{value}
	}
}
