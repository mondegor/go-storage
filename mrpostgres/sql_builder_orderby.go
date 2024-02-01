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
	SqlBuilderOrderBy struct {
		defaultOrderBy string
	}
)

func NewSqlBuilderOrderBy(ctx context.Context, defaultField string, defaultDirection mrenum.SortDirection) *SqlBuilderOrderBy {
	var defaultOrderBy string

	if defaultField != "" {
		defaultOrderBy = defaultField + " " + defaultDirection.String()
	} else {
		mrlog.Ctx(ctx).Warn().Caller(1).Msg("default sorting is not set")
	}

	return &SqlBuilderOrderBy{
		defaultOrderBy: defaultOrderBy,
	}
}

func NewSqlBuilderOrderByWithDefaultSort(ctx context.Context, defaultSort mrtype.SortParams) *SqlBuilderOrderBy {
	return NewSqlBuilderOrderBy(
		ctx,
		defaultSort.FieldName,
		defaultSort.Direction,
	)
}

func (b *SqlBuilderOrderBy) WrapWithDefault(field mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	if field != nil {
		return field
	}

	if b.defaultOrderBy == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return b.defaultOrderBy, []any{}
	}
}

func (b *SqlBuilderOrderBy) Join(fields ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
	fields = mrstorage.SqlBuilderPartFuncRemoveNil(fields)

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

func (b *SqlBuilderOrderBy) Field(name string, direction mrenum.SortDirection) mrstorage.SqlBuilderPartFunc {
	if name == "" {
		return nil
	}

	return func(paramNumber int) (string, []any) {
		return name + " " + direction.String(), []any{}
	}
}
