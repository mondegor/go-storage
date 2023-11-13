package mrsql

import (
	"fmt"
	"reflect"
	"strings"

	"github.com/mondegor/go-webcore/mrcore"
	"github.com/mondegor/go-webcore/mrenum"
	"github.com/mondegor/go-webcore/mrtype"
)

const (
	ModelNameEntityMetaOrderBy = "EntityMetaOrderBy"
)

type (
	EntityMetaOrderBy struct {
		fieldMap	map[string]bool
		defaultSort mrtype.SortParams
	}
)

// NewEntityMetaOrderBy - WARNING: use only when starting the main process
func NewEntityMetaOrderBy(entity any) (*EntityMetaOrderBy, error)  {
	rvt := reflect.TypeOf(entity)

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.FactoryErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	debugInfo := fmt.Sprintf("[%s] %s:", ModelNameEntityMetaOrderBy, rvt.String())

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
			mrcore.LogWarn(err)
			continue
		}

		var extMessage string

		if isDefault {
			meta.defaultSort.FieldName = sortName
			meta.defaultSort.Direction = sortDirection
			extMessage = ", default"
		}

		meta.fieldMap[sortName] = true
		debugInfo = fmt.Sprintf(
			"%s\n- %s(%d) -> %s %s%s;",
			debugInfo,
			rvt.Field(i).Name,
			i,
			sortName,
			sortDirection.String(),
			extMessage,
		)
	}

	mrcore.LogDebug(debugInfo)

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
		return errFunc("default field already exist")
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
