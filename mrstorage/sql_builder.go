package mrstorage

import (
	"github.com/mondegor/go-sysmess/mrtype"
	"github.com/mondegor/go-sysmess/mrtype/sortdirection"
)

type (
	// SQLBuilder - строитель SQL-условий для формирования частей запроса.
	// Предоставляет методы для построения SET, WHERE, ORDER BY и LIMIT конструкций.
	SQLBuilder interface {
		// Set - возвращает строитель для конструкции UPDATE SET.
		Set() SQLSetBuilder

		// Condition - возвращает строитель для конструкций WHERE и JOIN.
		Condition() SQLConditionBuilder

		// OrderBy - возвращает строитель для конструкции ORDER BY.
		OrderBy() SQLOrderByBuilder

		// Limit - возвращает строитель для конструкции LIMIT/OFFSET.
		Limit() SQLLimitBuilder
	}

	// SQLPart - параметризованная часть SQL-запроса.
	// Позволяет формировать SQL с настраиваемым префиксом и нумерацией аргументов.
	SQLPart interface {
		// WithPrefix - возвращает копию с указанным префиксом (например: " WHERE ").
		WithPrefix(sql string) SQLPart

		// WithStartArg - возвращает копию с указанным номером начального аргумента.
		WithStartArg(number int) SQLPart

		// Empty - проверяет, пуста ли часть SQL (отсутствует функция генерации).
		Empty() bool

		// ToSQL - возвращает SQL-выражение и список аргументов.
		ToSQL() (sql string, args []any)
	}

	// SQLSetBuilder - строитель для конструкции UPDATE SET.
	// Формирует список присваиваний вида "field1 = $1, field2 = $2".
	SQLSetBuilder interface {
		// Build - создаёт часть SQL из одной функции.
		Build(part SQLPartFunc) SQLPart

		// BuildComma - создаёт часть SQL, объединяя несколько функций через запятую.
		BuildComma(parts ...SQLPartFunc) SQLPart

		// BuildEntity - создаёт часть SQL на основе метаданных сущности.
		BuildEntity(entity any, parts ...SQLPartFunc) (SQLPart, error)

		// BuildFunc - создаёт часть SQL через функцию-помощник.
		BuildFunc(fn func(s SQLSetHelper) SQLPartFunc) SQLPart
	}

	// SQLSetHelper - помощник для построения выражений в конструкции SET.
	SQLSetHelper interface {
		// JoinComma - объединяет несколько выражений через запятую.
		JoinComma(fields ...SQLPartFunc) SQLPartFunc

		// Field - создаёт выражение присваивания одного поля.
		Field(name string, value any) SQLPartFunc

		// Fields - создаёт выражения присваивания нескольких полей.
		Fields(names []string, args []any) SQLPartFunc
	}

	// SQLConditionBuilder - строитель для конструкций WHERE и JOIN.
	// Формирует условия соединения и фильтрации записей.
	SQLConditionBuilder interface {
		// Build - создаёт часть SQL из одной функции.
		Build(part SQLPartFunc) SQLPart

		// BuildAnd - создаёт часть SQL, объединяя несколько условий через AND.
		BuildAnd(parts ...SQLPartFunc) SQLPart

		// BuildFunc - создаёт часть SQL через функцию-помощник.
		BuildFunc(fn func(c SQLConditionHelper) SQLPartFunc) SQLPart
	}

	// SQLConditionHelper - помощник для построения условий в WHERE и JOIN конструкциях.
	SQLConditionHelper interface {
		// JoinAnd - объединяет несколько условий через AND.
		JoinAnd(parts ...SQLPartFunc) SQLPartFunc

		// JoinOr - объединяет несколько условий через OR.
		JoinOr(parts ...SQLPartFunc) SQLPartFunc

		// Expr - создаёт произвольное SQL-выражение с параметрами.
		Expr(expr string, values ...any) SQLPartFunc

		// Equal - создаёт условие равенства (field = value).
		Equal(field string, value any) SQLPartFunc

		// NotEqual - создаёт условие неравенства (field <> value).
		NotEqual(field string, value any) SQLPartFunc

		// Less - создаёт условие меньше (field < value).
		Less(field string, value any) SQLPartFunc

		// LessOrEqual - создаёт условие меньше или равно (field <= value).
		LessOrEqual(field string, value any) SQLPartFunc

		// Greater - создаёт условие больше (field > value).
		Greater(field string, value any) SQLPartFunc

		// GreaterOrEqual - создаёт условие больше или равно (field >= value).
		GreaterOrEqual(field string, value any) SQLPartFunc

		// FilterEqual - создаёт условие равенства, пропуская пустые значения.
		FilterEqual(field string, value any) SQLPartFunc

		// FilterEqualString - создаёт условие равенства для строк, пропуская пустые.
		FilterEqualString(field, value string) SQLPartFunc

		// FilterEqualInt64 - создаёт условие равенства для int64, пропуская значение empty.
		FilterEqualInt64(field string, value, empty int64) SQLPartFunc

		// FilterEqualBool - создаёт условие равенства для bool, пропуская nil.
		FilterEqualBool(field string, value *bool) SQLPartFunc

		// FilterLike - создаёт условие LIKE для поиска подстроки.
		FilterLike(field, value string) SQLPartFunc

		// FilterLikePrefix - создаёт условие LIKE для поиска по префиксу.
		FilterLikePrefix(field, value string) SQLPartFunc

		// FilterLikeFields - создаёт условие LIKE для поиска по нескольким полям.
		FilterLikeFields(fields []string, value string) SQLPartFunc

		// FilterInArray - создаёт условие IN для массива JSON.
		// Параметр values поддерживает только слайсы, иначе возвращает nil.
		FilterInArray(jsonField string, values any) SQLPartFunc

		// FilterRangeInt64 - создаёт условие диапазона для int64.
		FilterRangeInt64(field string, value mrtype.RangeInt64, empty int64) SQLPartFunc

		// FilterRangeFloat64 - создаёт условие диапазона для float64.
		// Параметр qualityThreshold - порог точности для сравнения.
		FilterRangeFloat64(field string, value mrtype.RangeFloat64, empty, qualityThreshold float64) SQLPartFunc

		// FilterAnyOf - создаёт условие ANY() для проверки вхождения в массив.
		// Параметр values поддерживает только слайсы, иначе возвращает nil.
		FilterAnyOf(field string, values any) SQLPartFunc
	}

	// SQLOrderByBuilder - строитель для конструкции ORDER BY.
	// Формирует выражения сортировки результатов запроса.
	SQLOrderByBuilder interface {
		// Build - создаёт часть SQL из одной функции.
		Build(part SQLPartFunc) SQLPart

		// BuildComma - создаёт часть SQL, объединяя несколько полей сортировки через запятую.
		BuildComma(parts ...SQLPartFunc) SQLPart

		// BuildFunc - создаёт часть SQL через функцию-помощник.
		BuildFunc(fn func(o SQLOrderByHelper) SQLPartFunc) SQLPart
	}

	// SQLOrderByHelper - помощник для построения выражений сортировки в ORDER BY.
	SQLOrderByHelper interface {
		// JoinComma - объединяет несколько полей сортировки через запятую.
		JoinComma(fields ...SQLPartFunc) SQLPartFunc

		// Field - создаёт выражение сортировки по одному полю с указанием направления.
		Field(name string, direction sortdirection.Enum) SQLPartFunc
	}

	// SQLLimitBuilder - строитель для конструкции LIMIT/OFFSET.
	// Формирует выражения ограничения количества возвращаемых записей.
	SQLLimitBuilder interface {
		// Build - создаёт часть SQL для LIMIT с указанием смещения и размера.
		// index - номер начальной записи (OFFSET), size - количество записей (LIMIT).
		Build(index, size int) SQLPart
	}

	// SQLPartFunc - динамическая часть SQL-запроса, вычисляемая отложенно.
	// Функция вызывается тогда, когда известен номер первого параметра,
	// используемого в этой части запроса. Это позволяет корректно нумеровать
	// параметры ($1, $2, ...) при объединении нескольких частей SQL.
	SQLPartFunc func(argumentNumber int) (sql string, args []any)
)
