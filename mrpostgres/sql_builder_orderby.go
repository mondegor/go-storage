package mrpostgres

import (
	"context"
	"strings"

	"github.com/mondegor/go-storage/mrstorage"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

type (
	SQLBuilderOrderBy struct {
		defaultOrderBy string
	}
)

func NewSQLBuilderOrderBy(ctx context.Context, defaultSort mrtype.SortParams) *SQLBuilderOrderBy {
	var defaultOrderBy string

	if defaultSort.FieldName != "" {
		defaultOrderBy = defaultSort.FieldName + " " + defaultSort.Direction.String()
	} else {
		mrlog.Ctx(ctx).Warn().Caller(1).Msg("default sorting is not set")
	}

	return &SQLBuilderOrderBy{
		defaultOrderBy: defaultOrderBy,
	}
}

func (b *SQLBuilderOrderBy) DefaultField() mrstorage.SQLBuilderPartFunc {
	if b.defaultOrderBy == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return b.defaultOrderBy, []any{}
	}
}

func (b *SQLBuilderOrderBy) Join(fields ...mrstorage.SQLBuilderPartFunc) mrstorage.SQLBuilderPartFunc {
	fields = mrstorage.SQLBuilderPartFuncRemoveNil(fields)

	if len(fields) == 0 {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		var prepared []string

		for i := range fields {
			item, _ := fields[i](0)
			prepared = append(prepared, item)
		}

		return strings.Join(prepared, ", "), []any{}
	}
}

func (b *SQLBuilderOrderBy) Field(name string, direction mrenum.SortDirection) mrstorage.SQLBuilderPartFunc {
	if name == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " " + direction.String(), []any{}
	}
}
