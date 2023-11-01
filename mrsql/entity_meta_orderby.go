package mrsql

import (
    "fmt"
    "reflect"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    ModelNameEntityMetaOrderBy = "EntityMetaOrderBy"
)

type (
    EntityMetaOrderBy struct {
        FieldMap       map[string]string // field name -> db field name
        DefaultDbField string
    }
)

// NewEntityMetaOrderBy :WARNING: only for start service
func NewEntityMetaOrderBy(entity any) (*EntityMetaOrderBy, error)  {
    rvt := reflect.TypeOf(entity)

    for rvt.Kind() == reflect.Pointer {
        rvt = rvt.Elem()
    }

    if rvt.Kind() != reflect.Struct {
        return nil, mrcore.FactoryErrInternalInvalidType.Caller(1).New(rvt.Kind().String(), reflect.Struct.String())
    }

    debugInfo := fmt.Sprintf("[%s] %s:", ModelNameEntityMetaOrderBy, rvt.String())

    meta := EntityMetaOrderBy{
        FieldMap: make(map[string]string, 0),
    }

    for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
        fieldType := rvt.Field(i)
        sort := fieldType.Tag.Get(fieldTagSortByField)
        name := fieldType.Tag.Get(fieldTagJson)
        dbName := fieldType.Tag.Get(fieldTagDbFieldName)

        if sort != "" && sort != "on" && sort != "default" {
            mrcore.LogWarn(
                fmt.Sprintf(
                    "[%s] %s::%s sort = %s, expected value 'default' or 'on'",
                    ModelNameEntityMetaOrderBy,
                    rvt.String(),
                    rvt.Field(i).Name,
                    sort,
                ),
            )

            continue
        }

        if name == "" || dbName == "" {
            continue
        }

        var currentFieldDefault string

        if sort == "default" {
            if meta.DefaultDbField == "" {
                meta.DefaultDbField = dbName
                currentFieldDefault = "(default)"
            } else {
                mrcore.LogWarn(
                    fmt.Sprintf(
                        "[%s] %s::%s duplicate sort value 'default'",
                        ModelNameEntityMetaOrderBy,
                        rvt.String(),
                        rvt.Field(i).Name,
                    ),
                )
            }
        }

        meta.FieldMap[name] = dbName
        debugInfo = fmt.Sprintf("%s\n- %s(%d) -> %s%s;", debugInfo, name, i, dbName, currentFieldDefault)
    }

    mrcore.LogDebug(debugInfo)

    return &meta, nil
}
