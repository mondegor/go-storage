package txisolevel

// Уровни изоляции транзакции.
const (
	ReadUncommitted Enum = iota + 1 // read uncommitted
	ReadCommitted                   // read committed
	RepeatableRead                  // repeatable read
	Serializable                    // serializable
)

type (
	// Enum - перечисление элементов.
	Enum uint8
)

var enumKeys = map[Enum]string{ //nolint:gochecknoglobals
	ReadUncommitted: "READ_UNCOMMITTED",
	ReadCommitted:   "READ_COMMITTED",
	RepeatableRead:  "REPEATABLE_READ",
	Serializable:    "SERIALIZABLE",
}

// String - возвращает значение в виде строки.
func (e Enum) String() string {
	return enumKeys[e]
}
