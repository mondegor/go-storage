package mrentity

import (
    "reflect"
    "time"

    "github.com/mondegor/go-sysmess/mrerr"
)

const tagNameDb = "db"

func FilledFieldsToUpdate(entity any) (map[string]any, error) {
    rv := reflect.ValueOf(entity)

    for rv.Kind() == reflect.Pointer {
        rv = rv.Elem()
    }

    if rv.Kind() != reflect.Struct {
        return nil, mrerr.FactoryInternalInvalidType.Caller(1).New(rv.Kind().String(), reflect.Struct.String())
    }

    if !rv.IsValid() {
        return nil, mrerr.FactoryInternalInvalidData.Caller(1).New(rv)
    }

    values := make(map[string]any, 4)
    rvt := rv.Type()

    for i, cnt := 0, rv.NumField(); i < cnt; i++ {
        fieldType := rvt.Field(i)
        dbName := fieldType.Tag.Get(tagNameDb)

        if dbName == "" {
            continue
        }

        field := rv.Field(i)

        if !field.IsValid() {
            continue
        }

        switch fieldType.Type.Kind() {
        case reflect.String:
            if field.String() == "" {
                continue
            }

        case reflect.Int32, reflect.Int64:
            if field.Int() == 0 {
                continue
            }

        case reflect.Struct:
            v := field.Interface()
            value, ok := v.(time.Time)

            if ok && value.IsZero() {
                continue
            }

        default:
            continue
        }

        values[dbName] = field.Interface()
    }

    if len(values) == 0 {
        return values, FactoryInternalListOfFieldsIsEmpty.Caller(1).New()
    }

    return values, nil
}

