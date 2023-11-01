package mrpostgres

import (
    "fmt"
    "strings"

    "github.com/mondegor/go-storage/mrstorage"
)

type (
    SqlBuilderSet struct {
        fieldMap map[string]string // field name -> db field name
    }
)

func NewSqlBuilderSet() *SqlBuilderSet {
    return &SqlBuilderSet{}
}

func NewSqlBuilderSetWithFieldMap(fieldMap map[string]string) *SqlBuilderSet {
    return &SqlBuilderSet{
        fieldMap: fieldMap,
    }
}

func (b *SqlBuilderSet) DbName(name string) string {
    if b.fieldMap != nil {
        if dbName, ok := b.fieldMap[name]; ok {
            return dbName
        }
    }

    return name
}

func (b *SqlBuilderSet) Join(fields ...mrstorage.SqlBuilderPartFunc) mrstorage.SqlBuilderPartFunc {
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

func (b *SqlBuilderSet) Field(dbName string, value any) mrstorage.SqlBuilderPartFunc {
    return func (paramNumber int) (string, []any) {
        return fmt.Sprintf("%s = $%d", dbName, paramNumber), []any{value}
    }
}

func (b *SqlBuilderSet) Fields(dbNames []string, args []any) mrstorage.SqlBuilderPartFunc {
    return func (paramNumber int) (string, []any) {
        set := make([]string, len(dbNames))

        for i := range dbNames {
            set[i] = fmt.Sprintf("%s = $%d", dbNames[i], paramNumber)
            paramNumber++
        }

        return fmt.Sprintf("%s", strings.Join(set, ", ")), args
    }
}
