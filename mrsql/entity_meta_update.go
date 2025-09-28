package mrsql

import (
	"context"
	"fmt"
	"reflect"
	"regexp"
	"strings"
	"time"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
)

const (
	// ModelNameEntityMetaUpdate - название сущности.
	ModelNameEntityMetaUpdate = "EntityMetaUpdate"

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

var regexpDbName = regexp.MustCompile(`^[a-zA-Z_][a-zA-Z0-9_]*$`)

// NewEntityMetaUpdate - создаёт объект EntityMetaUpdate.
func NewEntityMetaUpdate(logger mrlog.Logger, entity any) (*EntityMetaUpdate, error) {
	rvt := reflect.TypeOf(entity)
	logger = logger.WithAttrs("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaUpdate, rvt.String()))

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mr.ErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	var debugInfo []string

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
			logger.Warn(context.Background(), "parse tag update warning, skipped", "error", err)

			continue
		}

		fieldType := field.Type
		isPointer := false

		if fieldType.Kind() == reflect.Pointer {
			fieldType = fieldType.Elem()
			isPointer = true
		}

		if !checkEntityMetaUpdateFieldType(fieldType) {
			logger.Warn(
				context.Background(),
				"check field type warning, skipped",
				"error", fmt.Errorf("field %s of type %s is not supported", rvt.Field(i).Name, fieldType.Kind()),
			)

			continue
		}

		meta.fieldsInfo[i] = fieldInfo{
			kind:      fieldType.Kind(),
			isPointer: isPointer,
			dbName:    dbName,
		}

		if logger.Enabled(mrlog.LevelDebug) {
			debugInfo = append(
				debugInfo,
				fmt.Sprintf(
					"- %s(%d, %s) -> %s;",
					rvt.Field(i).Name, i, rvt.Field(i).Type, dbName,
				),
			)
		}
	}

	if len(debugInfo) > 0 {
		logger.Debug(context.Background(), strings.Join(debugInfo, "\n"))
	}

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

// FieldsForUpdate - возвращает список полей и их значения для использования их при формировании SQL запроса.
func (m *EntityMetaUpdate) FieldsForUpdate(entity any) ([]string, []any, error) {
	rv := reflect.ValueOf(entity)

	for rv.Kind() == reflect.Pointer {
		rv = rv.Elem()
	}

	if !rv.IsValid() {
		return nil, nil, mr.ErrInternal.New("reflect.entity", rv)
	}

	if rv.Type().String() != m.structName {
		return nil, nil, mr.ErrInternalInvalidType.New(rv.Type().String(), m.structName)
	}

	fields := make([]string, 0, len(m.fieldsInfo))
	args := make([]any, 0, cap(fields))

	for i, info := range m.fieldsInfo {
		field := rv.Field(i)

		if !field.IsValid() {
			return nil, nil, mr.ErrInternal.New("reflect.field", field)
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

			value, ok := v.(time.Time)
			if !ok {
				return nil, nil, mr.ErrInternal.New("reflect.field.struct", field)
			}

			if !info.isPointer && value.IsZero() {
				continue
			}

		default:
			return nil, nil, mr.ErrInternal.New("reflect.field.undefined", field)
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
