package mrsql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrlog"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameEntityMetaOrderBy = "EntityMetaOrderBy"
)

type (
	EntityMetaOrderBy struct {
		fieldMap    map[string]bool
		defaultSort mrtype.SortParams
	}
)

// NewEntityMetaOrderBy - WARNING: use only when starting the main process
func NewEntityMetaOrderBy(ctx context.Context, entity any) (*EntityMetaOrderBy, error) {
	rvt := reflect.TypeOf(entity)
	logger := mrlog.Ctx(ctx).With().Str("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaOrderBy, rvt.String())).Logger()

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.FactoryErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	debugInfo := ""

	meta := EntityMetaOrderBy{
		fieldMap: make(map[string]bool, 0),
	}

	for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
		fieldType := rvt.Field(i)
		sort := fieldType.Tag.Get(fieldTagSortByField)

		if sort == "" {
			continue
		}

		sortName, isDefault, sortDirection, err := parseTagSort(rvt, sort, meta.defaultSort.FieldName == "")
		if err != nil {
			logger.Warn().Caller(1).Err(err).Msg("parse tag sort warning")
			continue
		}

		var extMessage string

		if isDefault {
			meta.defaultSort.FieldName = sortName
			meta.defaultSort.Direction = sortDirection
			extMessage = ", default"
		}

		meta.fieldMap[sortName] = true

		if logger.Level() <= mrlog.DebugLevel {
			debugInfo = fmt.Sprintf(
				"%s\n- %s(%d, %s) -> %s %s%s;",
				debugInfo,
				rvt.Field(i).Name,
				i,
				rvt.Field(i).Type,
				sortName,
				sortDirection.String(),
				extMessage,
			)
		}
	}

	logger.Debug().Msg(debugInfo)

	return &meta, nil
}

func (m *EntityMetaOrderBy) CheckField(name string) bool {
	_, ok := m.fieldMap[name]

	return ok
}

func (m *EntityMetaOrderBy) DefaultSort() mrtype.SortParams {
	return m.defaultSort
}

func parseTagSort(rvt reflect.Type, value string, canBeDefault bool) (string, bool, mrenum.SortDirection, error) {
	parsed := strings.Split(value, ",")
	count := len(parsed)

	errFunc := func(errString string) (string, bool, mrenum.SortDirection, error) {
		return "", false, 0, fmt.Errorf(
			"[%s] %s: parse error in '%s': %s",
			ModelNameEntityMetaOrderBy,
			rvt.String(),
			value,
			errString,
		)
	}

	if count > 3 {
		return errFunc("incorrect value")
	}

	if parsed[0] == "" {
		return errFunc("field name is required")
	}

	if !regexpDbName.MatchString(parsed[0]) {
		return errFunc("field name is incorrect")
	}

	isDefault := false

	if count > 1 {
		if parsed[1] != "default" {
			return errFunc("the second parameter can only be equal to 'default'")
		}

		isDefault = true
	}

	if !canBeDefault && isDefault {
		return errFunc("default field already exists")
	}

	sortDirection := mrenum.SortDirectionASC

	if count > 2 {
		parsed[2] = strings.ToUpper(parsed[2])

		if err := sortDirection.ParseAndSet(strings.ToUpper(parsed[2])); err != nil {
			return errFunc("the third parameter can only be equal to 'asc' or 'desc'")
		}
	}

	return parsed[0], isDefault, sortDirection, nil
}
