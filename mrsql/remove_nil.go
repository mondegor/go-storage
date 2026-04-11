package mrsql

import (
	"github.com/mondegor/go-storage/mrstorage"
)

// SQLPartFuncRemoveNil - удаляет все nil элементы из переданного списка функций SQLPartFunc.
// Используется для удаления необязательных частей SQL-запроса, которые не были инициализированы.
func SQLPartFuncRemoveNil(parts []mrstorage.SQLPartFunc) []mrstorage.SQLPartFunc {
	for i := 0; i < len(parts); i++ {
		if parts[i] != nil {
			continue
		}

		parts2 := parts[i:]
		for i2 := 1; i2 < len(parts2); i2++ {
			if parts2[i2] != nil {
				parts[i] = parts2[i2]
				i++
			}
		}

		clear(parts[i:]) // zero/nil out the obsolete elements, for GC

		return parts[:i]
	}

	return parts
}
