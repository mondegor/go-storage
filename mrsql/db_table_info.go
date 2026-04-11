package mrsql

type (
	// DBTableInfo - информация о таблице БД.
	DBTableInfo struct {
		Name       string // Name - имя таблицы в БД
		PrimaryKey string // PrimaryKey - имя первичного ключа таблицы
	}
)
