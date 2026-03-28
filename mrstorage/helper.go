package mrstorage

// ToSQL - comment func.
func ToSQL(part SQLPart) string {
	if part == nil {
		return ""
	}

	sql, _ := part.ToSQL()

	return sql
}
