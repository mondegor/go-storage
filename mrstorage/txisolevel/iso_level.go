package txisolevel

// Уровни изоляции транзакций в БД.
// Определяют степень видимости изменений, внесённых другими транзакциями.
const (
	ReadUncommitted Enum = iota + 1 // read uncommitted - самый низкий уровень, возможны "грязные" чтения
	ReadCommitted                   // read committed - чтение только зафиксированных данных (по умолчанию)
	RepeatableRead                  // repeatable read - гарантия повторного чтения без изменений
	Serializable                    // serializable - полная изоляция, последовательное выполнение
)

type (
	// Enum - перечисление уровней изоляции транзакций.
	Enum uint8
)

var enumKeys = map[Enum]string{ //nolint:gochecknoglobals
	ReadUncommitted: "READ_UNCOMMITTED",
	ReadCommitted:   "READ_COMMITTED",
	RepeatableRead:  "REPEATABLE_READ",
	Serializable:    "SERIALIZABLE",
}

// String - возвращает строковое представление уровня изоляции.
func (e Enum) String() string {
	if v, ok := enumKeys[e]; ok {
		return v
	}

	return "UNKNOWN"
}
