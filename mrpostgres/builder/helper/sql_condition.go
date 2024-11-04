package helper

import (
	"fmt"
	"reflect"
	"strconv"
	"strings"

	"github.com/mondegor/go-webcore/mrlib"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrsql"
	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLCondition - объект для создания независимой части SQL используемой в условий (WHERE, JOIN).
	SQLCondition struct{}
)

// NewSQLCondition - создаёт объект SQLCondition.
func NewSQLCondition() *SQLCondition {
	return &SQLCondition{}
}

// JoinAnd - возвращает указанные SQL условия соединённые оператором AND.
func (h *SQLCondition) JoinAnd(parts ...mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	return h.join(" AND ", parts)
}

// JoinOr - возвращает указанные SQL условия соединённые оператором OR.
func (h *SQLCondition) JoinOr(parts ...mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	return h.join(" OR ", parts)
}

// Expr - возвращает простое условие, например: "field_name BETWEEN 1000 AND 2000".
// Но если выражение пустое, то возвращается nil.
func (h *SQLCondition) Expr(expr string) mrstorage.SQLPartFunc {
	if expr == "" {
		return nil
	}

	return func(_ int) (string, []any) {
		return expr, nil
	}
}

// ExprWithValue - возвращает условие с аргументом, например: "UPPER(field_name) = %s".
func (h *SQLCondition) ExprWithValue(expr string, value any) mrstorage.SQLPartFunc {
	return func(argumentNumber int) (string, []any) {
		return fmt.Sprintf(expr, "$"+strconv.Itoa(argumentNumber)), []any{value}
	}
}

// Equal - возвращает строгое условие равенства.
func (h *SQLCondition) Equal(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, "=")
}

// NotEqual - возвращает строгое условие неравенства.
func (h *SQLCondition) NotEqual(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, "<>")
}

// Less - возвращает строгое условие меньше.
func (h *SQLCondition) Less(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, "<")
}

// LessOrEqual - возвращает строгое условие меньше равно.
func (h *SQLCondition) LessOrEqual(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, "<=")
}

// Greater - возвращает строгое условие больше.
func (h *SQLCondition) Greater(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, ">")
}

// GreaterOrEqual - возвращает строгое условие больше или равно.
func (h *SQLCondition) GreaterOrEqual(field string, value any) mrstorage.SQLPartFunc {
	return h.makeCompare(field, value, ">=")
}

// FilterEqual - возвращает условие равенства UUID если значение не пустое, иначе возвращается nil.
func (h *SQLCondition) FilterEqual(field string, value any) mrstorage.SQLPartFunc {
	rv := reflect.ValueOf(value)

	if !rv.IsValid() || rv.IsZero() {
		return nil
	}

	return h.makeCompare(field, value, "=")
}

// FilterEqualString - возвращает условие равенства строки если значение не пустое, иначе возвращается nil.
func (h *SQLCondition) FilterEqualString(field, value string) mrstorage.SQLPartFunc {
	if value == "" {
		return nil
	}

	return h.makeCompare(field, value, "=")
}

// FilterEqualInt64 - возвращает условие равенства целого числа если значение не пустое, иначе возвращается nil.
func (h *SQLCondition) FilterEqualInt64(field string, value, empty int64) mrstorage.SQLPartFunc {
	if value == empty {
		return nil
	}

	return h.makeCompare(field, value, "=")
}

// FilterEqualBool - возвращает условие равенства bool если значение не nil, иначе возвращается nil.
func (h *SQLCondition) FilterEqualBool(field string, value *bool) mrstorage.SQLPartFunc {
	if value == nil {
		return nil
	}

	return h.makeCompare(field, *value, "=")
}

// FilterLike - возвращает условие LIKE если значение не пустое, иначе возвращается nil.
func (h *SQLCondition) FilterLike(field, value string) mrstorage.SQLPartFunc {
	return h.FilterLikeFields([]string{field}, value)
}

// FilterLikeFields - возвращает условие LIKE для нескольких полей если значение не пустое, иначе возвращается nil.
func (h *SQLCondition) FilterLikeFields(fields []string, value string) mrstorage.SQLPartFunc {
	if value == "" {
		return nil
	}

	// sample: (field_name LIKE '%%' || $1 || '%%' OR ...)
	return func(argumentNumber int) (string, []any) {
		var buf strings.Builder

		buf.Grow(30 * len(fields))
		buf.WriteByte('(')

		for i := range fields {
			if i > 0 {
				buf.WriteString(" OR ")
			}

			buf.WriteString(fields[i])
			buf.WriteString(" LIKE '%%' || $")
			buf.WriteString(strconv.Itoa(argumentNumber))
			buf.WriteString(" || '%%'")
		}

		buf.WriteByte(')')

		return buf.String(), []any{value}
	}
}

// FilterRangeInt64 - возвращает интервальное условие для целых чисел если значения Min, Max не пустые, иначе возвращается nil.
func (h *SQLCondition) FilterRangeInt64(field string, value mrtype.RangeInt64, empty int64) mrstorage.SQLPartFunc {
	if value.Min != empty {
		if value.Max != empty {
			if value.Min > value.Max {
				return nil
			}

			return func(argumentNumber int) (string, []any) {
				return "(" + field + " BETWEEN $" + strconv.Itoa(argumentNumber) + " AND $" + strconv.Itoa(argumentNumber+1) + ")", []any{value.Min, value.Max}
			}
		}

		return h.makeCompare(field, value.Min, ">=")
	} else if value.Max != empty {
		return h.makeCompare(field, value.Max, "<=")
	}

	return nil
}

// FilterRangeFloat64 - возвращает интервальное условие для вещественных чисел если значения Min, Max не пустые, иначе возвращается nil.
func (h *SQLCondition) FilterRangeFloat64(field string, value mrtype.RangeFloat64, empty, qualityThreshold float64) mrstorage.SQLPartFunc {
	if !mrlib.EqualFloat(value.Min, empty, qualityThreshold) {
		if !mrlib.EqualFloat(value.Max, empty, qualityThreshold) {
			if value.Min > value.Max {
				return nil
			}

			return func(argumentNumber int) (string, []any) {
				return "(" + field + " BETWEEN $" + strconv.Itoa(argumentNumber) + " AND $" + strconv.Itoa(argumentNumber+1) + ")",
					[]any{value.Min - qualityThreshold, value.Max + qualityThreshold}
			}
		}

		return h.makeCompare(field, value.Min-qualityThreshold, ">=")
	} else if value.Max != empty {
		return h.makeCompare(field, value.Max+qualityThreshold, "<=")
	}

	return nil
}

// FilterAnyOf - возвращает условие (= ANY), которое проверяет, чтобы хотя бы один элемент из списка был равен значению указанного поля.
// Параметр 'values' поддерживает только слайсы с хотя бы одним значением, иначе вернётся nil.
func (h *SQLCondition) FilterAnyOf(field string, values any) mrstorage.SQLPartFunc {
	rv := reflect.ValueOf(values)

	if rv.Kind() != reflect.Slice || rv.Len() == 0 {
		return nil
	}

	// sample: field_name = ANY($1)
	return func(argumentNumber int) (string, []any) {
		return field + " = ANY($" + strconv.Itoa(argumentNumber) + ")", []any{values}
	}
}

// FilterInOf - возвращает условие (IN), которое проверяет, чтобы хотя бы один элемент из списка был равен значению указанного поля.
// Параметр 'values' support only slices else the func returns nil.
func (h *SQLCondition) FilterInOf(field string, values any) mrstorage.SQLPartFunc {
	rv := reflect.ValueOf(values)

	if rv.Kind() != reflect.Slice || rv.Len() == 0 {
		return nil
	}

	args := make([]any, rv.Len())

	for i := 0; i < len(args); i++ {
		args[i] = rv.Index(i).Interface()
	}

	// sample: field_name IN($1, $2, ...)
	return func(argumentNumber int) (string, []any) {
		var buf strings.Builder

		buf.Grow(len(field) + 4 + 3*len(args)) // len(field) + " IN(" + "$N," * len(args) - 1
		buf.WriteString(field)
		buf.WriteString(" IN(")

		for i := range args {
			if i > 0 {
				buf.WriteByte(',')
			}

			buf.WriteByte('$')
			buf.WriteString(strconv.Itoa(argumentNumber + i))
		}

		buf.WriteByte(')')

		return buf.String(), args
	}
}

func (h *SQLCondition) join(separator string, conditions []mrstorage.SQLPartFunc) mrstorage.SQLPartFunc {
	conditions = mrsql.SQLPartFuncRemoveNil(conditions)

	if len(conditions) == 0 {
		return nil
	}

	// sample: (cond1 AND cond2 AND ...)
	return func(argumentNumber int) (string, []any) {
		var (
			buf  strings.Builder
			args []any
		)

		buf.WriteByte('(')

		for i := range conditions {
			if i > 0 {
				buf.WriteString(separator)
			}

			item, itemArgs := conditions[i](argumentNumber + len(args))
			buf.WriteString(item)

			args = mrsql.MergeArgs(args, itemArgs)
		}

		buf.WriteByte(')')

		return buf.String(), args
	}
}

func (h *SQLCondition) makeCompare(field string, value any, sign string) mrstorage.SQLPartFunc {
	return func(argumentNumber int) (string, []any) {
		return field + " " + sign + " $" + strconv.Itoa(argumentNumber), []any{value}
	}
}
