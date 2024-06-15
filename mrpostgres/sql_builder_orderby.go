package mrpostgres

import (
	"context"
	"errors"
	"strings"

	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"

	"github.com/mondegor/go-storage/mrstorage"
)

type (
	// SQLBuilderOrderBy - comment struct.
	SQLBuilderOrderBy struct {
		defaultOrderBy string
	}
)

// NewSQLBuilderOrderBy - comment func.
func NewSQLBuilderOrderBy(ctx context.Context, defaultSort mrtype.SortParams) *SQLBuilderOrderBy {
	var defaultOrderBy string

	if defaultSort.FieldName != "" {
		defaultOrderBy = defaultSort.FieldName + " " + defaultSort.Direction.String()
	} else {
		mrlog.Ctx(ctx).Warn().Err(errors.New("default sorting is not set")).Send()
	}

	return &SQLBuilderOrderBy{
		defaultOrderBy: defaultOrderBy,
	}
}

// DefaultField - comment method.
func (b *SQLBuilderOrderBy) DefaultField() mrstorage.SQLBuilderPartFunc {
	if b.defaultOrderBy == "" {
		return nil
	}

	return func(_ int) (string, []any) {
		return b.defaultOrderBy, nil
	}
}

// Join - comment method.
func (b *SQLBuilderOrderBy) Join(fields ...mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	fields = mrstorage.SQLBuilderPartFuncRemoveNil(fields)

	if len(fields) == 0 {
		return nil
	}

	return func(_ int) (string, []any) {
		var prepared []string

		for i := range fields {
			item, _ := fields[i](0)
			prepared = append(prepared, item)
		}

		return strings.Join(prepared, ", "), nil
	}
}

// Field - comment method.
func (b *SQLBuilderOrderBy) Field(name string, direction mrenum.SortDirection) mrstorage.SQLBuilderPartFunc {
	if name == "" {
		return nil
	}

	return func(_ int) (string, []any) {
		return name + " " + direction.String(), nil
	}
}
