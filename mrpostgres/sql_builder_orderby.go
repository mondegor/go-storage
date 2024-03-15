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

func NewSqlBuilderOrderBy(ctx context.Context, defaultSort mrtype.SortParams) *SqlBuilderOrderBy {
	var defaultOrderBy string

	if defaultSort.FieldName != "" {
		defaultOrderBy = defaultSort.FieldName + " " + defaultSort.Direction.String()
	} else {
		mrlog.Ctx(ctx).Warn().Caller(1).Msg("default sorting is not set")
	}

	return &SqlBuilderOrderBy{
		defaultOrderBy: defaultOrderBy,
	}
}

func (b *SqlBuilderOrderBy) DefaultField() mrstorage.SqlBuilderPartFunc {
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
