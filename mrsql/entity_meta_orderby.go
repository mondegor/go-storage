package mrsql

import (
	"context"
	"fmt"
	"reflect"
	"strings"

	"github.com/mondegor/go-sysmess/mrerr/mr"
	"github.com/mondegor/go-sysmess/mrlog"
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/enums"
)

const (
	// ModelNameEntityMetaOrderBy - название сущности.
	ModelNameEntityMetaOrderBy = "EntityMetaOrderBy"

	fieldTagSortByField = "sort"
)

type (
	// EntityMetaOrderBy - объект для управления порядком следования записей БД.
	// Информация о порядке следования считывается из тегов структуры.
	EntityMetaOrderBy struct {
		fieldMap    map[string]bool
		defaultSort mrtype.SortParams
	}

	parsedTagSort struct {
		SortName      string
		IsDefault     bool
		SortDirection enums.SortDirection
	}
)

// NewEntityMetaOrderBy - создаёт объект EntityMetaOrderBy.
func NewEntityMetaOrderBy(logger mrlog.Logger, entity any) (*EntityMetaOrderBy, error) {
	rvt := reflect.TypeOf(entity)
	logger = logger.WithAttrs("object", fmt.Sprintf("[%s] %s", ModelNameEntityMetaOrderBy, rvt.String()))

	for rvt.Kind() == reflect.Pointer {
		rvt = rvt.Elem()
	}

	if rvt.Kind() != reflect.Struct {
		return nil, mr.ErrInternalInvalidType.New(rvt.Kind().String(), reflect.Struct.String())
	}

	var debugInfo []string

	meta := EntityMetaOrderBy{
		fieldMap: make(map[string]bool),
	}

	for i, cnt := 0, rvt.NumField(); i < cnt; i++ {
		fieldType := rvt.Field(i)
		sort := fieldType.Tag.Get(fieldTagSortByField)

		if sort == "" {
			continue
		}

		parsed, err := parseTagSort(rvt, sort, meta.defaultSort.FieldName == "")
		if err != nil {
			logger.Warn(context.Background(), "parse tag sort warning, skipped")

			continue
		}

		var extMessage string

		if parsed.IsDefault {
			meta.defaultSort.FieldName = parsed.SortName
			meta.defaultSort.Direction = parsed.SortDirection
			extMessage = ", default"
		}

		meta.fieldMap[parsed.SortName] = true

		if logger.Enabled(mrlog.LevelDebug) {
			debugInfo = append(
				debugInfo,
				fmt.Sprintf(
					"- %s(%d, %s) -> %s %s%s;",
					rvt.Field(i).Name, i, rvt.Field(i).Type, parsed.SortName, parsed.SortDirection.String(), extMessage,
				),
			)
		}
	}

	if len(debugInfo) > 0 {
		logger.Debug(context.Background(), strings.Join(debugInfo, "\n"))
	}

	return &meta, nil
}

// CheckField - сообщает, зарегистрировано ли указанное поле в распарсенной структуре.
func (m *EntityMetaOrderBy) CheckField(name string) bool {
	_, ok := m.fieldMap[name]

	return ok
}

// DefaultSort - возвращает данные о сортировке по умолчанию.
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

	tagSort := parsedTagSort{
		SortName:      parsed[0],
		IsDefault:     isDefault,
		SortDirection: enums.SortDirectionASC,
	}

	if count > 2 {
		sortDirection, err := enums.ParseSortDirection(strings.ToUpper(parsed[2]))
		if err != nil {
			return errFunc("the third parameter can only be equal to 'asc' or 'desc'")
		}

		tagSort.SortDirection = sortDirection
	}

	return tagSort, nil
}
