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
	ModelNameEntityMetaOrderBy = "EntityMetaOrderBy" // ModelNameEntityMetaOrderBy - название сущности

	fieldTagSortByField = "sort"
)

type (
	// EntityMetaOrderBy - comment struct.
	EntityMetaOrderBy struct {
		fieldMap    map[string]bool
		defaultSort mrtype.SortParams
	}

	parsedTagSort struct {
		SortName      string
		IsDefault     bool
		SortDirection mrenum.SortDirection
	}
)

// NewEntityMetaOrderBy - создаёт объект EntityMetaOrderBy.
// WARNING: use only when starting the main process.
func NewEntityMetaOrderBy(ctx context.Context, entity any) (*EntityMetaOrderBy, error) {
	rvt := reflect.TypeOf(entity)
	logger := mrlog.Ctx(ctx).With().Str("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaOrderBy, rvt.String())).Logger()

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mrcore.ErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
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

		parsed, err := parseTagSort(rvt, sort, meta.defaultSort.FieldName == "")
		if err != nil {
			logger.Warn().Err(err).Msg("parse tag sort warning, skipped")

			continue
		}

		var extMessage string

		if parsed.IsDefault {
			meta.defaultSort.FieldName = parsed.SortName
			meta.defaultSort.Direction = parsed.SortDirection
			extMessage = ", default"
		}

		meta.fieldMap[parsed.SortName] = true

		if logger.Level() <= mrlog.DebugLevel {
			debugInfo = fmt.Sprintf(
				"%s\n- %s(%d, %s) -> %s %s%s;",
				debugInfo,
				rvt.Field(i).Name,
				i,
				rvt.Field(i).Type,
				parsed.SortName,
				parsed.SortDirection.String(),
				extMessage,
			)
		}
	}

	logger.Debug().Msg(debugInfo)

	return &meta, nil
}

// CheckField - comment method.
func (m *EntityMetaOrderBy) CheckField(name string) bool {
	_, ok := m.fieldMap[name]

	return ok
}

// DefaultSort - comment method.
func (m *EntityMetaOrderBy) DefaultSort() mrtype.SortParams {
	return m.defaultSort
}

func parseTagSort(rvt reflect.Type, value string, canBeDefault bool) (parsedTagSort, error) {
	parsed := strings.Split(value, ",")
	count := len(parsed)

	errFunc := func(errString string) (parsedTagSort, error) {
		return parsedTagSort{}, fmt.Errorf(
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

	return parsedTagSort{
		SortName:      parsed[0],
		IsDefault:     isDefault,
		SortDirection: sortDirection,
	}, nil
}
