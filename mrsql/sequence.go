package mrsql

// SequenceName - возвращает название последовательности используемой для получения ID.
func SequenceName(table DBTableInfo) string {
	return table.Name + "_" + table.PrimaryKey + "_seq"
}
