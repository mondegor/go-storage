package mrsql

import (
	"context"
	"fmt"
	"reflect"
	"time"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrlog"
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
func NewEntityMetaUpdate(ctx context.Context, entity any) (*EntityMetaUpdate, error) {
	rvt := reflect.TypeOf(entity)
	logger := mrlog.Ctx(ctx).With().Str("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaUpdate, rvt.String())).Logger()

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.FactoryErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	debugInfo := ""

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
			logger.Warn().Caller(1).Err(err).Msg("parse tag update warning")
			continue
		}

		meta.fieldsInfo[i] = fieldInfo{
			kind:   fieldType.Type.Kind(),
			dbName: dbName,
		}

		if logger.Level() <= mrlog.DebugLevel {
			debugInfo = fmt.Sprintf("%s\n- %s(%d) -> %s;", debugInfo, rvt.Field(i).Name, i, dbName)
		}
	}

	logger.Debug().Msg(debugInfo)

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
		return nil, nil, mrcore.FactoryErrInternal.WithAttr("reflect.entity", rv).New()
	}

	var fields []string
	var args []any

	for i, info := range m.fieldsInfo {
		field := rv.Field(i)

		if !field.IsValid() {
			continue
		}

		switch info.kind {
		case reflect.String, reflect.Slice:
			if field.Len() == 0 {
				continue
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64:
			if field.Int() == 0 {
				continue
			}

		case reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64:
			if field.Uint() == 0 {
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
