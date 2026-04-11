package mrsql

// SequenceName - возвращает имя PostgreSQL-последовательности для получения ID.
// Формирует имя по шаблону: {TableName}_{PrimaryKey}_seq
// Например: для таблицы "users" с ключом "id" -> "users_id_seq".
func SequenceName(table DBTableInfo) string {
	return table.Name + "_" + table.PrimaryKey + "_seq"
}
