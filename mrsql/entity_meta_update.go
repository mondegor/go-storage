package mrsql

import (
    "fmt"
    "reflect"
    "time"

    "github.com/mondegor/go-webcore/mrcore"
)

const (
    ModelNameEntityMetaUpdate = "EntityMetaUpdate"
)

type (
    EntityMetaUpdate struct {
        StructName string
        Fields map[string]string // field name -> db field name
        FieldsInfo map[int]fieldInfo // field index -> fieldInfo
    }

    fieldInfo struct {
        kind reflect.Kind
        dbName string
    }
)

// NewEntityMetaUpdate :WARNING: only for start service
func NewEntityMetaUpdate(entity any) (*EntityMetaUpdate, error) {
    rvt := reflect.TypeOf(entity)

    for rvt.Kind() == reflect.Pointer {
        rvt = rvt.Elem()
    }

    if rvt.Kind() != reflect.Struct {
        return nil, mrcore.FactoryErrInternalInvalidType.Caller(1).New(rvt.Kind().String(), reflect.Struct.String())
    }

    debugInfo := fmt.Sprintf("[%s] %s:", ModelNameEntityMetaUpdate, rvt.String())

    meta := EntityMetaUpdate{
        StructName: rvt.String(),
        Fields: make(map[string]string, 0),
        FieldsInfo: make(map[int]fieldInfo, 0),
    }

    for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
        fieldType := rvt.Field(i)
        update := fieldType.Tag.Get(fieldTagFreeUpdate)
        name := fieldType.Tag.Get(fieldTagJson)
        dbName := fieldType.Tag.Get(fieldTagDbFieldName)

        if update != "" && update != "on" {
            mrcore.LogWarn(
                fmt.Sprintf(
                    "[%s] %s::%s update = %s, expected value = 'on'",
                    ModelNameEntityMetaUpdate,
                    rvt.String(),
                    rvt.Field(i).Name,
                    update,
                ),
            )

            continue
        }

        if name == "" || dbName == "" {
            continue
        }

        meta.Fields[name] = dbName
        meta.FieldsInfo[i] = fieldInfo{
            kind: fieldType.Type.Kind(),
            dbName: dbName,
        }

        debugInfo = fmt.Sprintf("%s\n- %s(%d) -> %s;", debugInfo, name, i, dbName)
    }

    mrcore.LogDebug(debugInfo)

    return &meta, nil
}

func FieldsForUpdate(meta *EntityMetaUpdate, entity any) ([]string, []any, error) {
    rv := reflect.ValueOf(entity)

    for rv.Kind() == reflect.Pointer {
        rv = rv.Elem()
    }

    if rv.Type().String() != meta.StructName {
        return nil, nil, mrcore.FactoryErrInternalInvalidType.Caller(1).New(rv.Type().String(), meta.StructName)
    }

    if !rv.IsValid() {
        return nil, nil, mrcore.FactoryErrInternalInvalidData.Caller(1).New(rv)
    }

    var fields []string
    var args []any

    for i, info := range meta.FieldsInfo {
        field := rv.Field(i)

        if !field.IsValid() {
            continue
        }

        switch info.kind {
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

        fields = append(fields, info.dbName)
        args = append(args, field.Interface())
    }

    return fields, args, nil
}

