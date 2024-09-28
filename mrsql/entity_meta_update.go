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
	ModelNameEntityMetaUpdate = "EntityMetaUpdate" // ModelNameEntityMetaUpdate - название сущности

	fieldTagDBFieldName = "db"
	fieldTagFieldUpdate = "upd"
)

type (
	// EntityMetaUpdate - объект для управления динамическим обновлением записей в БД.
	// Информация об обновлении считывается из тегов структуры.
	EntityMetaUpdate struct {
		structName string
		fieldsInfo map[int]fieldInfo // field index -> fieldInfo
	}

	fieldInfo struct {
		kind      reflect.Kind
		isPointer bool
		dbName    string
	}
)

// NewEntityMetaUpdate - создаёт объект EntityMetaUpdate.
// WARNING: use only when starting the main process.
func NewEntityMetaUpdate(ctx context.Context, entity any) (*EntityMetaUpdate, error) {
	rvt := reflect.TypeOf(entity)
	logger := mrlog.Ctx(ctx).With().Str("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaUpdate, rvt.String())).Logger()

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.ErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	debugInfo := ""

	meta := EntityMetaUpdate{
		structName: rvt.String(),
		fieldsInfo: make(map[int]fieldInfo),
	}

	for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
		field := rvt.Field(i)
		update := field.Tag.Get(fieldTagFieldUpdate)
		dbName := field.Tag.Get(fieldTagDBFieldName)

		if update == "" {
			continue
		}

		dbName, err := parseTagUpdate(rvt, update, dbName)
		if err != nil {
			logger.Warn().Err(err).Msg("parse tag update warning, skipped")

			continue
		}

		fieldType := field.Type
		isPointer := false

		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
			isPointer = true
		}

		if !checkEntityMetaUpdateFieldType(fieldType) {
			logger.Warn().Err(
				fmt.Errorf("field %s of type %s is not supported", rvt.Field(i).Name, fieldType.Kind()),
			).Msg("check field type warning, skipped")

			continue
		}

		meta.fieldsInfo[i] = fieldInfo{
			kind:      fieldType.Kind(),
			isPointer: isPointer,
			dbName:    dbName,
		}

		if logger.Level() <= mrlog.DebugLevel {
			debugInfo = fmt.Sprintf("%s\n- %s(%d, %s) -> %s;", debugInfo, rvt.Field(i).Name, i, rvt.Field(i).Type, dbName)
		}
	}

	logger.Debug().Msg(debugInfo)

	return &meta, nil
}

func parseTagUpdate(rvt reflect.Type, value, dbName string) (string, error) {
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

// FieldsForUpdate - comment method.
func (m *EntityMetaUpdate) FieldsForUpdate(entity any) ([]string, []any, error) {
	rv := reflect.ValueOf(entity)

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if rv.Type().String() != m.structName {
		return nil, nil, mrcore.ErrInternalInvalidType.New(rv.Type().String(), m.structName)
	}

	if !rv.IsValid() {
		return nil, nil, mrcore.ErrInternal.New().WithAttr("reflect.entity", rv)
	}

	fields := make([]string, 0, len(m.fieldsInfo))
	args := make([]any, 0, cap(fields))

	for i, info := range m.fieldsInfo {
		field := rv.Field(i)

		if !field.IsValid() {
			return nil, nil, mrcore.ErrInternal.New().WithAttr("reflect.field", field)
		}

		if info.isPointer {
			if field.IsNil() {
				continue
			}

			field = rv.Field(i).Elem()
		}

		switch info.kind {
		case reflect.String, reflect.Slice: // empty slice === nil
			if !info.isPointer && field.Len() == 0 {
				continue
			}

		case reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
			reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
			reflect.Float32, reflect.Float64,
			reflect.Bool, reflect.Array:
			if !info.isPointer && field.IsZero() {
				continue
			}

		case reflect.Struct:
			v := field.Interface()

			if value, ok := v.(time.Time); ok {
				if !info.isPointer && value.IsZero() {
					continue
				}
			} else {
				return nil, nil, mrcore.ErrInternal.New().WithAttr("reflect.field.struct", field)
			}

		default:
			return nil, nil, mrcore.ErrInternal.New().WithAttr("reflect.field.undefined", field)
		}

		fields = append(fields, info.dbName)
		args = append(args, field.Interface())
	}

	return fields, args, nil
}

func checkEntityMetaUpdateFieldType(fieldType reflect.Type) bool {
	switch fieldType.Kind() {
	case reflect.String,
		reflect.Int, reflect.Int8, reflect.Int16, reflect.Int32, reflect.Int64,
		reflect.Uint, reflect.Uint8, reflect.Uint16, reflect.Uint32, reflect.Uint64,
		reflect.Float32, reflect.Float64,
		reflect.Bool:
		return true

	case reflect.Array:
		return fieldType.String() == "uuid.UUID"

	case reflect.Slice:
		return fieldType.Elem().Name() == "uint8" // byte

	case reflect.Struct:
		return fieldType.String() == "time.Time"
	}

	return false
}
