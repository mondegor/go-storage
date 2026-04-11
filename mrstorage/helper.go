package mrstorage

// ToSQL - преобразует часть SQL-запроса в строку.
// Если part равен nil, возвращает пустую строку.
// Игнорирует аргументы, возвращает только SQL-выражение.
func ToSQL(part SQLPart) string {
	if part == nil {
		return ""
	}

	sql, _ := part.ToSQL()

	return sql
}
