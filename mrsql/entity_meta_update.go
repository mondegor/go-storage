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
		structName string
		fieldsInfo map[int]fieldInfo // field index -> fieldInfo
	}

	fieldInfo struct {
		kind   reflect.Kind
		dbName string
	}
)

// NewEntityMetaUpdate - WARNING: use only when starting the main process
func NewEntityMetaUpdate(entity any) (*EntityMetaUpdate, error) {
	rvt := reflect.TypeOf(entity)

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.FactoryErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	debugInfo := fmt.Sprintf("[%s] %s:", ModelNameEntityMetaUpdate, rvt.String())

	meta := EntityMetaUpdate{
		structName: rvt.String(),
		fieldsInfo: make(map[int]fieldInfo, 0),
	}

	for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
		fieldType := rvt.Field(i)
		update := fieldType.Tag.Get(fieldTagFieldUpdate)
		dbName := fieldType.Tag.Get(fieldTagDBFieldName)

		if update == "" {
			continue
		}

		dbName, err := parseTagUpdate(rvt, update, dbName)

		if err != nil {
			mrcore.LogWarn(err)
			continue
		}

		meta.fieldsInfo[i] = fieldInfo{
			kind:   fieldType.Type.Kind(),
			dbName: dbName,
		}

		debugInfo = fmt.Sprintf("%s\n- %s(%d) -> %s;", debugInfo, rvt.Field(i).Name, i, dbName)
	}

	mrcore.LogDebug(debugInfo)

	return &meta, nil
}

func parseTagUpdate(rvt reflect.Type, value string, dbName string) (string, error) {
	errFunc := func(errString string) (string, error) {
		return "", fmt.Errorf(
			"[%s] %s: parse error in '%s': %s",
			ModelNameEntityMetaUpdate,
			rvt.String(),
			value,
			errString,
		)
	}

	if value == "+" {
		if dbName == "" {
			return errFunc(fmt.Sprintf("tag '%s' is empty", fieldTagDBFieldName))
		}

		if !regexpDbName.MatchString(dbName) {
			return errFunc(fmt.Sprintf("value '%s' from '%s' is incorrect", dbName, fieldTagDBFieldName))
		}
	} else if dbName == "" {
		if !regexpDbName.MatchString(value) {
			return errFunc("value is incorrect")
		}

		dbName = value
	}

	return dbName, nil
}

func (m *EntityMetaUpdate) FieldsForUpdate(entity any) ([]string, []any, error) {
	rv := reflect.ValueOf(entity)

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Type().String() != m.structName {
		return nil, nil, mrcore.FactoryErrInternalInvalidType.New(rv.Type().String(), m.structName)
	}

	if !rv.IsValid() {
		return nil, nil, mrcore.FactoryErrInternalInvalidData.New(rv)
	}

	var fields []string
	var args []any

	for i, info := range m.fieldsInfo {
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
