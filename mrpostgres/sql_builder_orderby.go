package mrpostgres

import (
    "fmt"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
    "github.com/mondegor/go-webcore/mrcore"
    "github.com/mondegor/go-webcore/mrenum"
    "github.com/mondegor/go-webcore/mrtype"
)

type (
    SqlBuilderOrderBy struct {
        defaultOrderBy string
    }
)

func NewSqlBuilderOrderBy(defaultField string, defaultDirection mrenum.SortDirection) *SqlBuilderOrderBy {
    var defaultOrderBy string

    if defaultField != "" {
        defaultOrderBy = fmt.Sprintf("%s %s", defaultField, defaultDirection.String())
    } else {
        mrcore.LogWarning("default sorting is not set")
    }

    return &SqlBuilderOrderBy{
        defaultOrderBy: defaultOrderBy,
    }
}

func NewSqlBuilderOrderByWithDefaultSort(defaultSort mrtype.SortParams) *SqlBuilderOrderBy {
    return NewSqlBuilderOrderBy(
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

        return fmt.Sprintf("%s", strings.Join(prepared, ", ")), []any{}
    }
}

func (b *SqlBuilderOrderBy) Field(name string, direction mrenum.SortDirection) mrstorage.SqlBuilderPartFunc {
    if name == "" {
        return nil
    }

    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s %s", name, direction.String()), []any{}
    }
}
