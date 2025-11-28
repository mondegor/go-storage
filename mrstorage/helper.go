package mrstorage

// ToSQL - comment func.
func ToSQL(part SQLPart) string {
	sql, _ := part.ToSQL()

	return sql
}
