package mrpostgres

import (
    "fmt"
    "strings"

    "github.com/mondegor/go-storage/mrentity"
    "github.com/mondegor/go-storage/mrstorage"
)

type (
    SqlBuilderOrderBy struct {
        defaultOrderBy string
        fieldMap map[string]string // field name -> db field name
    }
)

func NewSqlBuilderOrderBy(defaultDbField string, defaultDirection mrentity.SortDirection) *SqlBuilderOrderBy {
    return NewSqlBuilderOrderByWithFieldMap(nil, defaultDbField, defaultDirection)
}

func NewSqlBuilderOrderByWithFieldMap(fieldMap map[string]string, defaultDbField string, defaultDirection mrentity.SortDirection) *SqlBuilderOrderBy {
    var defaultOrderBy string

    if defaultDbField != "" {
        defaultOrderBy = fmt.Sprintf("%s %s", defaultDbField, defaultDirection.String())
    }

    return &SqlBuilderOrderBy{
        fieldMap: fieldMap,
        defaultOrderBy: defaultOrderBy,
    }
}

func (b *SqlBuilderOrderBy) DbName(name string) string {
    if b.fieldMap != nil {
        if dbName, ok := b.fieldMap[name]; ok {
            return dbName
        }
    }

    return name
}

func (b *SqlBuilderOrderBy) Join(fields ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
    fields = mrstorage.SqlBuilderPartFuncRemoveNil(fields)

    if len(fields) == 0 {
        if b.defaultOrderBy == "" {
            return nil
        }

        return func(paramNumber int) (string, []any) {
            return b.defaultOrderBy, []any{}
        }
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

func (b *SqlBuilderOrderBy) Field(dbName string, direction mrentity.SortDirection) mrstorage.SqlBuilderPartFunc {
    if dbName == "" {
        return nil
    }

    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s %s", dbName, direction.String()), []any{}
    }
}
