package mrstorage

type (
	// MetaGetter - интерфейс для предоставления информации о таблице БД.
	MetaGetter interface {
		TableName() string
		PrimaryName() string
		Condition() SQLBuilderPart
	}
)
